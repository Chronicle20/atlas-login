package handler

import (
	"atlas-login/account"
	as "atlas-login/account/session"
	"atlas-login/configuration"
	"atlas-login/kafka/producer"
	"atlas-login/session"
	"atlas-login/socket/writer"
	"context"
	"github.com/Chronicle20/atlas-model/model"
	"github.com/Chronicle20/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const LoginHandle = "LoginHandle"

type LoginRequest struct {
	name           string
	password       string
	hwid           []byte
	gameRoomClient uint32
	gameStartMode  byte
	unknown1       byte
}

func (l *LoginRequest) Name() string {
	return l.name
}

func (l *LoginRequest) Password() string {
	return l.password
}

func (l *LoginRequest) GameRoomClient() uint32 {
	return l.gameRoomClient
}

func (l *LoginRequest) GameStartMode() byte {
	return l.gameStartMode
}

func ReadLoginRequest(reader *request.Reader) *LoginRequest {
	name := reader.ReadAsciiString()
	password := reader.ReadAsciiString()
	hwid := reader.ReadBytes(16)
	gameRoomClient := reader.ReadUint32()
	gameStartMode := reader.ReadByte()
	unknown1 := reader.ReadByte()

	return &LoginRequest{
		name:           name,
		password:       password,
		hwid:           hwid,
		gameRoomClient: gameRoomClient,
		gameStartMode:  gameStartMode,
		unknown1:       unknown1,
	}
}

func LoginHandleFunc(l logrus.FieldLogger, ctx context.Context, wp writer.Producer) func(s session.Model, r *request.Reader) {
	authTemporaryBanFunc := session.Announce(l)(wp)(writer.AuthTemporaryBan)
	authPermanentBanFunc := session.Announce(l)(wp)(writer.AuthPermanentBan)

	return func(s session.Model, r *request.Reader) {
		p := ReadLoginRequest(r)
		l.Debugf("Reading [%s] message. body={name=%s, password=%s, gameRoomClient=%d, gameStartMode=%d}", LoginHandle, p.Name(), p.Password(), p.GameRoomClient(), p.GameStartMode())

		resp, err := as.CreateLogin(l, ctx)(s.SessionId(), s.AccountId(), p.Name(), p.Password(), "")
		if err != nil {
			announceError(l, ctx, wp)(s, writer.SystemError1)
			return
		}

		t := s.Tenant()
		if resp.Code == "OK" || resp.Code == writer.LicenseAgreement {
			var a account.Model
			a, err = account.GetByName(l, ctx)(p.Name())
			if err != nil {
				announceError(l, ctx, wp)(s, writer.SystemError1)
				return
			}
			s = session.SetAccountId(a.Id())(t.Id(), s.SessionId())
			session.SessionCreated(producer.ProviderImpl(l)(ctx), t)(s)

			if resp.Code == "OK" {
				err = issueSuccess(l, s, wp)(a)
				if err != nil {
					l.WithError(err).Errorf("Unable to issue success to account.")
					return
				}
				if t.Region() == "JMS" {
					issueServerInformation(l, ctx, wp)(s)
				}
				return
			}

		}

		if resp.Code != writer.Banned {
			announceError(l, ctx, wp)(s, resp.Code)
			return
		}

		if resp.Until != 0 {
			err = authTemporaryBanFunc(s, writer.AuthTemporaryBanBody(l, t)(resp.Until, resp.Reason))
			if err != nil {
				l.WithError(err).Errorf("Unable to show account is temporary banned.")
			}
			return
		}

		err = authPermanentBanFunc(s, writer.AuthPermanentBanBody(l, t))
		if err != nil {
			l.WithError(err).Errorf("Unable to show account is permanently banned.")
		}
	}
}

func issueSuccess(l logrus.FieldLogger, s session.Model, wp writer.Producer) model.Operator[account.Model] {
	authSuccessFunc := session.Announce(l)(wp)(writer.AuthSuccess)
	return func(a account.Model) error {
		c, err := configuration.GetConfiguration()
		if err != nil {
			l.WithError(err).Errorf("Unable to get configuration.")
			return err
		}
		t := s.Tenant()
		sc, err := c.FindServer(t.Id().String())
		if err != nil {
			l.WithError(err).Errorf("Unable to find server configuration.")
			return err
		}

		err = authSuccessFunc(s, writer.AuthSuccessBody(t)(a.Id(), a.Name(), a.Gender(), sc.UsesPIN, a.PIC()))
		if err != nil {
			l.WithError(err).Errorf("Unable to show successful authorization for account %d", a.Id())
		}
		return err
	}
}

func announceError(l logrus.FieldLogger, _ context.Context, wp writer.Producer) func(s session.Model, reason string) {
	authLoginFailedFunc := session.Announce(l)(wp)(writer.AuthLoginFailed)
	return func(s session.Model, reason string) {
		err := authLoginFailedFunc(s, writer.AuthLoginFailedBody(l, s.Tenant())(reason))
		if err != nil {
			l.WithError(err).Errorf("Unable to issue [%s].", writer.AuthLoginFailed)
		}
	}
}
