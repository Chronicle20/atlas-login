package handler

import (
	as "atlas-login/account/session"
	"atlas-login/kafka/producer"
	"atlas-login/session"
	"atlas-login/socket/model"
	"atlas-login/socket/writer"
	"atlas-login/world/channel"
	"context"
	"github.com/Chronicle20/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const CharacterSelectedHandle = "CharacterSelectedHandle"

func CharacterSelectedHandleFunc(l logrus.FieldLogger, ctx context.Context, wp writer.Producer) func(s session.Model, r *request.Reader) {
	serverIpFunc := session.Announce(l)(wp)(writer.ServerIP)
	return func(s session.Model, r *request.Reader) {
		characterId := r.ReadUint32()
		var sMacAddressWithHDDSerial = ""
		var sMacAddressWithHDDSerial2 = ""

		t := s.Tenant()
		if t.Region() == "GMS" {
			if t.MajorVersion() > 12 {
				sMacAddressWithHDDSerial = r.ReadAsciiString()
				sMacAddressWithHDDSerial2 = r.ReadAsciiString()
			}
		}
		l.Debugf("Character [%d] selected for login to channel [%d:%d]. hwid [%s] hwid [%s].", characterId, s.WorldId(), s.ChannelId(), sMacAddressWithHDDSerial, sMacAddressWithHDDSerial2)

		c, err := channel.GetById(l, ctx)(s.WorldId(), s.ChannelId())
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve channel information being logged in to.")
			err = serverIpFunc(s, writer.ServerIPBodySimpleError(l)(writer.ServerIPCodeServerUnderInspection))
			if err != nil {
				l.WithError(err).Errorf("Unable to write server ip response due to error.")
				return
			}
			return
		}

		err = as.UpdateState(l, producer.ProviderImpl(l)(ctx))(s.SessionId(), s.AccountId(), 2, model.ChannelSelect{IPAddress: c.IpAddress(), Port: uint16(c.Port()), CharacterId: characterId})
		if err != nil {
			return
		}

		//resp, err := as.UpdateState(l, ctx)(s.SessionId(), s.AccountId(), 2)
		//if err != nil || resp.Code != "OK" {
		//	l.WithError(err).Errorf("Unable to update session for character [%d] attempting to login.", characterId)
		//	err = serverIpFunc(s, writer.ServerIPBodySimpleError(l)(writer.ServerIPCodeTooManyConnectionRequests))
		//	if err != nil {
		//		l.WithError(err).Errorf("Unable to write server ip response due to error.")
		//		return
		//	}
		//	return
		//}
		//
		//err = serverIpFunc(s, writer.ServerIPBody(l, t)(c.IpAddress(), uint16(c.Port()), characterId))
		//if err != nil {
		//	l.WithError(err).Errorf("Unable to write server ip response due to error.")
		//	return
		//}
	}
}
