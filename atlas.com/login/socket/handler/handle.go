package handler

import (
	"atlas-login/session"
	"atlas-login/socket/writer"
	"atlas-login/tracing"
	"github.com/Chronicle20/atlas-socket/request"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type MessageValidator func(l logrus.FieldLogger, span opentracing.Span) func(s session.Model) bool

const NoOpValidator = "NoOpValidator"

func NoOpValidatorFunc(_ logrus.FieldLogger, _ opentracing.Span) func(_ session.Model) bool {
	return func(_ session.Model) bool {
		return true
	}
}

const LoggedInValidator = "LoggedInValidator"

func LoggedInValidatorFunc(_ logrus.FieldLogger, _ opentracing.Span) func(s session.Model) bool {
	return func(s session.Model) bool {
		//v := account.IsLoggedIn(l, span)(s.AccountId())
		//if !v {
		//	l.Errorf("Attempting to process a request when the account %d is not logged in.", s.SessionId())
		//}
		//return v
		return true
	}
}

type MessageHandler func(l logrus.FieldLogger, span opentracing.Span, wp writer.Producer) func(s session.Model, r *request.Reader)

const NoOpHandler = "NoOpHandler"

func NoOpHandlerFunc(_ logrus.FieldLogger, _ opentracing.Span, _ writer.Producer) func(_ session.Model, _ *request.Reader) {
	return func(_ session.Model, _ *request.Reader) {
	}
}

func AdaptHandler(l logrus.FieldLogger, name string, v MessageValidator, h MessageHandler, wp writer.Producer) request.Handler {
	return func(sessionId uuid.UUID, r request.Reader) {
		fl := l.WithField("session", sessionId.String())
		sl, span := tracing.StartSpan(fl, name)

		s, ok := session.GetRegistry().Get(sessionId)
		if !ok {
			sl.Errorf("Unable to locate session %d", sessionId)
			return
		}

		if v(sl, span)(s) {
			h(sl, span, wp)(s, &r)
			s = session.UpdateLastRequest()(s.SessionId())
		}
		span.Finish()
	}
}
