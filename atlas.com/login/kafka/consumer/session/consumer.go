package session

import (
	"atlas-login/account"
	"atlas-login/configuration"
	consumer2 "atlas-login/kafka/consumer"
	"atlas-login/kafka/producer"
	"atlas-login/session"
	model2 "atlas-login/socket/model"
	"atlas-login/socket/writer"
	"atlas-login/world"
	"context"
	"github.com/Chronicle20/atlas-kafka/consumer"
	"github.com/Chronicle20/atlas-kafka/handler"
	"github.com/Chronicle20/atlas-kafka/message"
	"github.com/Chronicle20/atlas-kafka/topic"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/Chronicle20/atlas-tenant"
	"github.com/sirupsen/logrus"
	"sort"
)

func InitConsumers(l logrus.FieldLogger) func(func(config consumer.Config, decorators ...model.Decorator[consumer.Config])) func(consumerGroupId string) {
	return func(rf func(config consumer.Config, decorators ...model.Decorator[consumer.Config])) func(consumerGroupId string) {
		return func(consumerGroupId string) {
			rf(consumer2.NewConfig(l)("account_session_status_event")(EnvEventStatusTopic)(consumerGroupId), consumer.SetHeaderParsers(consumer.SpanHeaderParser, consumer.TenantHeaderParser))
		}
	}
}

func InitHandlers(l logrus.FieldLogger) func(tenant tenant.Model) func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
	return func(tenant tenant.Model) func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
		return func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
			return func(rf func(topic string, handler handler.Handler) (string, error)) {
				var t string
				t, _ = topic.EnvProvider(l)(EnvEventStatusTopic)()
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleCreatedAccountSessionStatusEvent(tenant, wp))))
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleLicenseAgreementAccountSessionStatusEvent(tenant, wp))))
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleErrorAccountSessionStatusEvent(tenant, wp))))
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleStateChangedAccountSessionStatusEvent(tenant, wp))))
			}
		}
	}
}
func handleCreatedAccountSessionStatusEvent(t tenant.Model, wp writer.Producer) func(l logrus.FieldLogger, ctx context.Context, e statusEvent[createdStatusEventBody]) {
	return func(l logrus.FieldLogger, ctx context.Context, e statusEvent[createdStatusEventBody]) {
		if e.Type != EventStatusTypeCreated {
			return
		}

		if !t.Is(tenant.MustFromContext(ctx)) {
			return
		}

		session.IfPresentById(t)(e.SessionId, func(s session.Model) error {
			a, err := account.GetById(l, ctx)(e.AccountId)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve account [%d] that had a session [%s] created for it.", e.AccountId, s.SessionId().String())
				return err
			}

			s = session.SetAccountId(a.Id())(t.Id(), s.SessionId())
			session.SessionCreated(producer.ProviderImpl(l)(ctx), t)(s)

			c, err := configuration.Get()
			if err != nil {
				l.WithError(err).Errorf("Unable to get configuration.")
				return err
			}

			sc, err := c.FindServer(t.Id())
			if err != nil {
				l.WithError(err).Errorf("Unable to find server configuration.")
				return err
			}

			err = session.Announce(l)(wp)(writer.AuthSuccess)(s, writer.AuthSuccessBody(t)(a.Id(), a.Name(), a.Gender(), sc.UsesPIN, a.PIC()))
			if err != nil {
				l.WithError(err).Errorf("Unable to show successful authorization for account %d", a.Id())
				return err
			}

			if t.Region() == "JMS" {
				_ = announceServerInformation(l)(ctx)(wp)(s)
			}

			return err
		})
	}
}

func handleLicenseAgreementAccountSessionStatusEvent(t tenant.Model, wp writer.Producer) message.Handler[statusEvent[any]] {
	return func(l logrus.FieldLogger, ctx context.Context, e statusEvent[any]) {
		if e.Type != EventStatusTypeRequestLicenseAgreement {
			return
		}

		if !t.Is(tenant.MustFromContext(ctx)) {
			return
		}

		session.IfPresentById(t)(e.SessionId, func(s session.Model) error {
			a, err := account.GetById(l, ctx)(e.AccountId)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve account [%d] that had a session [%s] created for it.", e.AccountId, s.SessionId().String())
				return err
			}

			s = session.SetAccountId(a.Id())(t.Id(), s.SessionId())
			session.SessionCreated(producer.ProviderImpl(l)(ctx), t)(s)

			return announceError(l)(ctx)(wp)("LICENSE_AGREEMENT")(s)
		})
	}
}

func handleErrorAccountSessionStatusEvent(t tenant.Model, wp writer.Producer) func(l logrus.FieldLogger, ctx context.Context, e statusEvent[errorStatusEventBody]) {
	return func(l logrus.FieldLogger, ctx context.Context, e statusEvent[errorStatusEventBody]) {
		if e.Type != EventStatusTypeError {
			return
		}

		if !t.Is(tenant.MustFromContext(ctx)) {
			return
		}

		if e.Body.Code != EventStatusErrorCodeDeletedOrBlocked {
			session.IfPresentById(t)(e.SessionId, announceError(l)(ctx)(wp)(e.Body.Code))
		} else if e.Body.Until != 0 {
			session.IfPresentById(t)(e.SessionId, announceTemporaryBan(l)(ctx)(wp)(e.Body.Until, e.Body.Reason))
		} else {
			session.IfPresentById(t)(e.SessionId, announcePermanentBan(l)(ctx)(wp))
		}
	}
}

func announcePermanentBan(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) model.Operator[session.Model] {
		t := tenant.MustFromContext(ctx)
		return func(wp writer.Producer) model.Operator[session.Model] {
			authPermanentBanFunc := session.Announce(l)(wp)(writer.AuthPermanentBan)
			return func(s session.Model) error {
				err := authPermanentBanFunc(s, writer.AuthPermanentBanBody(l, t))
				if err != nil {
					l.WithError(err).Errorf("Unable to show account is permanently banned.")
				}
				return err
			}
		}
	}
}

func announceTemporaryBan(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(until uint64, reason byte) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(until uint64, reason byte) model.Operator[session.Model] {
		t := tenant.MustFromContext(ctx)
		return func(wp writer.Producer) func(until uint64, reason byte) model.Operator[session.Model] {
			authTemporaryBanFunc := session.Announce(l)(wp)(writer.AuthTemporaryBan)
			return func(until uint64, reason byte) model.Operator[session.Model] {
				return func(s session.Model) error {
					err := authTemporaryBanFunc(s, writer.AuthTemporaryBanBody(l, t)(until, reason))
					if err != nil {
						l.WithError(err).Errorf("Unable to show account is temporary banned.")
					}
					return err
				}
			}
		}
	}
}

func announceError(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(reason string) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(reason string) model.Operator[session.Model] {
		return func(wp writer.Producer) func(reason string) model.Operator[session.Model] {
			authLoginFailedFunc := session.Announce(l)(wp)(writer.AuthLoginFailed)
			return func(reason string) model.Operator[session.Model] {
				return func(s session.Model) error {
					err := authLoginFailedFunc(s, writer.AuthLoginFailedBody(l, s.Tenant())(reason))
					if err != nil {
						l.WithError(err).Errorf("Unable to issue [%s].", writer.AuthLoginFailed)
						return err
					}
					return nil
				}
			}
		}
	}
}

func announceServerInformation(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) model.Operator[session.Model] {
		ws, err := world.GetAll(l, ctx, world.ChannelLoadDecorator(l, ctx))
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve worlds to display to session.")
		}
		sort.Slice(ws, func(i, j int) bool {
			return ws[i].Id() < ws[j].Id()
		})

		if len(ws) == 0 {
			l.Warnf("Responding with no valid worlds.")
		}

		return func(wp writer.Producer) model.Operator[session.Model] {
			return model.ThenOperator(announceServerList(l)(ctx)(wp)(ws), model.Operators[session.Model](announceLastWorld(l)(ctx)(wp), announceRecommendedWorlds(l)(ctx)(wp)(ws)))
		}
	}
}

func announceRecommendedWorlds(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(ws []world.Model) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(ws []world.Model) model.Operator[session.Model] {
		return func(wp writer.Producer) func(ws []world.Model) model.Operator[session.Model] {
			serverListRecommendationFunc := session.Announce(l)(wp)(writer.ServerListRecommendations)
			return func(ws []world.Model) model.Operator[session.Model] {
				return func(s session.Model) error {
					var rs = make([]world.Recommendation, 0)
					for _, x := range ws {
						if x.Recommended() {
							rs = append(rs, x.Recommendation())
						}
					}
					err := serverListRecommendationFunc(s, writer.ServerListRecommendationsBody(rs))
					if err != nil {
						l.WithError(err).Errorf("Unable to issue recommended worlds")
						return err
					}
					return nil
				}
			}
		}
	}
}

func announceLastWorld(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) model.Operator[session.Model] {
		return func(wp writer.Producer) model.Operator[session.Model] {
			selectWorldFunc := session.Announce(l)(wp)(writer.SelectWorld)
			return func(s session.Model) error {
				err := selectWorldFunc(s, writer.SelectWorldBody(0))
				if err != nil {
					l.WithError(err).Errorf("Unable to identify the last world a account was logged into")
					return err
				}
				return nil
			}
		}
	}
}

func announceServerList(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(ws []world.Model) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(ws []world.Model) model.Operator[session.Model] {
		return func(wp writer.Producer) func(ws []world.Model) model.Operator[session.Model] {
			serverListEntryFunc := session.Announce(l)(wp)(writer.ServerListEntry)
			serverListEndFunc := session.Announce(l)(wp)(writer.ServerListEnd)
			return func(ws []world.Model) model.Operator[session.Model] {
				return func(s session.Model) error {
					for _, x := range ws {
						err := serverListEntryFunc(s, writer.ServerListEntryBody(s.Tenant())(x.Id(), x.Name(), x.State(), x.EventMessage(), x.ChannelLoad()))
						if err != nil {
							l.WithError(err).Errorf("Unable to write server list entry for [%d]", x.Id())
						}
					}
					err := serverListEndFunc(s, writer.ServerListEndBody)
					if err != nil {
						l.WithError(err).Errorf("Unable to complete writing the server list")
						return err
					}
					return nil
				}
			}
		}
	}
}

func handleStateChangedAccountSessionStatusEvent(t tenant.Model, wp writer.Producer) message.Handler[statusEvent[stateChangedEventBody[model2.ChannelSelect]]] {
	return func(l logrus.FieldLogger, ctx context.Context, e statusEvent[stateChangedEventBody[model2.ChannelSelect]]) {
		if e.Type != EventStatusTypeStateChanged {
			return
		}

		if !t.Is(tenant.MustFromContext(ctx)) {
			return
		}

		session.IfPresentById(t)(e.SessionId, processStateReturn(l)(ctx)(wp)(e.AccountId, e.Body.State, e.Body.Params))
	}
}

func processStateReturn(l logrus.FieldLogger) func(ctx context.Context) func(wp writer.Producer) func(accountId uint32, state uint8, params model2.ChannelSelect) model.Operator[session.Model] {
	return func(ctx context.Context) func(wp writer.Producer) func(accountId uint32, state uint8, params model2.ChannelSelect) model.Operator[session.Model] {
		t := tenant.MustFromContext(ctx)
		return func(wp writer.Producer) func(accountId uint32, state uint8, params model2.ChannelSelect) model.Operator[session.Model] {
			serverIpFunc := session.Announce(l)(wp)(writer.ServerIP)
			return func(accountId uint32, state uint8, params model2.ChannelSelect) model.Operator[session.Model] {
				return func(s session.Model) error {
					if len(params.IPAddress) <= 0 {
						return nil
					}

					err := serverIpFunc(s, writer.ServerIPBody(l, t)(params.IPAddress, params.Port, params.CharacterId))
					if err != nil {
						l.WithError(err).Errorf("Unable to write server ip response due to error.")
						return err
					}
					return nil
				}
			}
		}
	}
}
