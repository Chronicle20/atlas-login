package handler

import (
	"atlas-login/session"
	"atlas-login/socket/writer"
	"github.com/Chronicle20/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"math/rand"
)

const CreateSecurityHandle = "CreateSecurityHandle"

func CreateSecurityHandleFunc(l logrus.FieldLogger, _ opentracing.Span, wp writer.Producer) func(s session.Model, r *request.Reader) {
	loginAuthFunc := session.Announce(l)(wp)(writer.LoginAuth)

	return func(s session.Model, _ *request.Reader) {
		loginScreen := [2]string{"MapLogin", "MapLogin1"}
		randomIndex := rand.Intn(len(loginScreen))

		err := loginAuthFunc(s, writer.LoginAuthBody(l)(loginScreen[randomIndex]))
		if err != nil {
			l.WithError(err).Errorf("Unable to announce login screen.")
		}
	}
}
