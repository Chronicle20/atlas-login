package session

import (
	"atlas-login/kafka/producer"
	"atlas-login/socket/writer"
	"atlas-login/tenant"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Announce(l logrus.FieldLogger) func(writerProducer writer.Producer) func(writerName string) func(s Model, bodyProducer writer.BodyProducer) error {
	return func(writerProducer writer.Producer) func(writerName string) func(s Model, bodyProducer writer.BodyProducer) error {
		return func(writerName string) func(s Model, bodyProducer writer.BodyProducer) error {
			return func(s Model, bodyProducer writer.BodyProducer) error {
				w, err := writerProducer(l, writerName)
				if err != nil {
					return err
				}

				if lock, ok := GetRegistry().GetLock(s.Tenant().Id, s.SessionId()); ok {
					lock.Lock()
					err = s.announceEncrypted(w(l)(bodyProducer))
					lock.Unlock()
					return err
				}
				return errors.New("invalid session")
			}
		}
	}
}

func SetAccountId(accountId uint32) func(tenantId uuid.UUID, id uuid.UUID) Model {
	return func(tenantId uuid.UUID, id uuid.UUID) Model {
		s := Model{}
		var ok bool
		if s, ok = GetRegistry().Get(tenantId, id); ok {
			s = s.setAccountId(accountId)
			GetRegistry().Update(s)
			return s
		}
		return s
	}
}

func UpdateLastRequest() func(tenantId uuid.UUID, id uuid.UUID) Model {
	return func(tenantId uuid.UUID, id uuid.UUID) Model {
		s := Model{}
		var ok bool
		if s, ok = GetRegistry().Get(tenantId, id); ok {
			s = s.updateLastRequest()
			GetRegistry().Update(s)
			return s
		}
		return s
	}
}

func SetWorldId(worldId byte) func(tenantId uuid.UUID, id uuid.UUID) Model {
	return func(tenantId uuid.UUID, id uuid.UUID) Model {
		s := Model{}
		var ok bool
		if s, ok = GetRegistry().Get(tenantId, id); ok {
			s = s.setWorldId(worldId)
			GetRegistry().Update(s)
			return s
		}
		return s
	}
}

func SetChannelId(channelId byte) func(tenantId uuid.UUID, id uuid.UUID) Model {
	return func(tenantId uuid.UUID, id uuid.UUID) Model {
		s := Model{}
		var ok bool
		if s, ok = GetRegistry().Get(tenantId, id); ok {
			s = s.setChannelId(channelId)
			GetRegistry().Update(s)
			return s
		}
		return s
	}
}

func SessionCreated(kp producer.Provider, tenant tenant.Model) func(s Model) {
	return func(s Model) {
		_ = kp(EnvEventTopicSessionStatus)(createdStatusEventProvider(tenant, s.SessionId(), s.AccountId()))
	}
}
