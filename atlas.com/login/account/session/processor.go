package session

import (
	"atlas-login/kafka/producer"
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Create(_ logrus.FieldLogger, kp producer.Provider) func(sessionId uuid.UUID, accountId uint32, accountName string, password string, ipAddress string) error {
	return func(sessionId uuid.UUID, accountId uint32, accountName string, password string, ipAddress string) error {
		return kp(EnvCommandTopic)(createCommandProvider(sessionId, accountId, accountName, password, ipAddress))
	}
}

func Destroy(l logrus.FieldLogger, kp producer.Provider) func(sessionId uuid.UUID, accountId uint32) {
	return func(sessionId uuid.UUID, accountId uint32) {
		l.Debugf("Destroying session for account [%d].", accountId)
		_ = kp(EnvCommandTopic)(logoutCommandProvider(sessionId, accountId))
	}
}

func UpdateState(l logrus.FieldLogger, ctx context.Context) func(sessionId uuid.UUID, accountId uint32, state int) (Model, error) {
	return func(sessionId uuid.UUID, accountId uint32, state int) (Model, error) {
		return updateState(l, ctx)(sessionId, accountId, state)
	}
}
