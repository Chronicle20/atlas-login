package account

import (
	consumer2 "atlas-login/kafka/consumer"
	"atlas-login/socket/writer"
	"context"
	"github.com/Chronicle20/atlas-kafka/consumer"
	"github.com/Chronicle20/atlas-kafka/handler"
	"github.com/Chronicle20/atlas-kafka/message"
	"github.com/Chronicle20/atlas-kafka/topic"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/Chronicle20/atlas-tenant"
	"github.com/sirupsen/logrus"
)

func InitConsumers(l logrus.FieldLogger) func(rf func(config consumer.Config, decorators ...model.Decorator[consumer.Config])) func(consumerGroupId string) {
	return func(rf func(config consumer.Config, decorators ...model.Decorator[consumer.Config])) func(consumerGroupId string) {
		return func(consumerGroupId string) {
			rf(consumer2.NewConfig(l)("account_status_event")(EnvEventTopicAccountStatus)(consumerGroupId), consumer.SetHeaderParsers(consumer.SpanHeaderParser, consumer.TenantHeaderParser))
		}
	}
}

func InitHandlers(l logrus.FieldLogger) func(tenant tenant.Model) func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
	return func(tenant tenant.Model) func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
		return func(wp writer.Producer) func(rf func(topic string, handler handler.Handler) (string, error)) {
			return func(rf func(topic string, handler handler.Handler) (string, error)) {
				t, _ := topic.EnvProvider(l)(EnvEventTopicAccountStatus)()
				_, _ = rf(t, message.AdaptHandler(message.PersistentConfig(handleAccountStatusEvent(tenant))))
			}
		}
	}
}

func handleAccountStatusEvent(ot tenant.Model) func(l logrus.FieldLogger, ctx context.Context, event statusEvent) {
	return func(l logrus.FieldLogger, ctx context.Context, event statusEvent) {
		t, err := tenant.FromContext(ctx)()
		if err != nil {
			l.WithError(err).Error("error getting tenant")
			return
		}

		if !t.Is(ot) {
			return
		}

		if event.Status == EventAccountStatusLoggedIn {
			getRegistry().Login(Key{Tenant: t, Id: event.AccountId})
		} else if event.Status == EventAccountStatusLoggedOut {
			getRegistry().Logout(Key{Tenant: t, Id: event.AccountId})
		}
	}
}
