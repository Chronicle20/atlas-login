package handler

import (
	"atlas-login/session"
	"atlas-login/socket/writer"
	"github.com/Chronicle20/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const CharacterViewAllPongHandle = "CharacterViewAllPongHandle"

func CharacterViewAllPongHandleFunc(l logrus.FieldLogger, _ opentracing.Span, _ writer.Producer) func(s session.Model, r *request.Reader) {
	return func(s session.Model, r *request.Reader) {
		var opt = r.ReadBool()
		var mode = "RESET"
		if opt {
			mode = "RENDER"
		}
		l.Debugf("View All Character PONG for account [%d]. mode [%s].", s.AccountId(), mode)
	}
}
