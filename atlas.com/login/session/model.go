package session

import (
	"github.com/Chronicle20/atlas-socket/crypto"
	"github.com/Chronicle20/atlas-tenant"
	"github.com/google/uuid"
	"math/rand"
	"net"
	"time"
)

type Model struct {
	id          uuid.UUID
	accountId   uint32
	worldId     byte
	channelId   byte
	con         net.Conn
	send        crypto.AESOFB
	recv        crypto.AESOFB
	encryptFunc crypto.EncryptFunc
	lastPacket  time.Time
	tenant      tenant.Model
	locale      byte
}

func NewSession(id uuid.UUID, t tenant.Model, locale byte, con net.Conn) Model {
	recvIv := []byte{byte(rand.Float64() * 255), byte(rand.Float64() * 255), byte(rand.Float64() * 255), byte(rand.Float64() * 255)}
	sendIv := []byte{byte(rand.Float64() * 255), byte(rand.Float64() * 255), byte(rand.Float64() * 255), byte(rand.Float64() * 255)}

	var send *crypto.AESOFB
	var recv *crypto.AESOFB
	if t.Region() == "GMS" && t.MajorVersion() <= 12 {
		send = crypto.NewAESOFB(sendIv, uint16(65535)-t.MajorVersion(), crypto.SetIvGenerator(crypto.FillIvZeroGenerator))
		recv = crypto.NewAESOFB(recvIv, t.MajorVersion(), crypto.SetIvGenerator(crypto.FillIvZeroGenerator))
	} else {
		send = crypto.NewAESOFB(sendIv, uint16(65535)-t.MajorVersion())
		recv = crypto.NewAESOFB(recvIv, t.MajorVersion())
	}

	hasMapleEncryption := true
	if t.Region() == "JMS" {
		hasMapleEncryption = false
	}

	return Model{
		id:          id,
		con:         con,
		send:        *send,
		recv:        *recv,
		encryptFunc: send.Encrypt(hasMapleEncryption, true),
		lastPacket:  time.Now(),
		tenant:      t,
		locale:      locale,
	}
}

func CloneSession(s Model) Model {
	return Model{
		id:          s.id,
		accountId:   s.accountId,
		worldId:     s.worldId,
		channelId:   s.channelId,
		con:         s.con,
		send:        s.send,
		recv:        s.recv,
		encryptFunc: s.encryptFunc,
		lastPacket:  s.lastPacket,
		tenant:      s.tenant,
		locale:      s.locale,
	}
}

func (s *Model) setAccountId(accountId uint32) Model {
	ns := CloneSession(*s)
	ns.accountId = accountId
	return ns
}

func (s *Model) SessionId() uuid.UUID {
	return s.id
}

func (s *Model) AccountId() uint32 {
	return s.accountId
}

func (s *Model) Tenant() tenant.Model {
	return s.tenant
}

func (s *Model) announceEncrypted(b []byte) error {
	tmp := make([]byte, len(b)+4)
	copy(tmp, b)
	tmp = append([]byte{0, 0, 0, 0}, b...)
	wip := s.encryptFunc(tmp)
	_, err := s.con.Write(wip)
	return err
}

func (s *Model) announce(b []byte) error {
	_, err := s.con.Write(b)
	return err
}

func (s *Model) WriteHello() error {
	return s.announce(WriteHello(nil)(s.tenant.MajorVersion(), s.tenant.MinorVersion(), s.send.IV(), s.recv.IV(), s.locale))
}

func (s *Model) ReceiveAESOFB() *crypto.AESOFB {
	return &s.recv
}

func (s *Model) GetRemoteAddress() net.Addr {
	return s.con.RemoteAddr()
}

func (s *Model) setWorldId(worldId byte) Model {
	ns := CloneSession(*s)
	ns.worldId = worldId
	return ns
}

func (s *Model) setChannelId(channelId byte) Model {
	ns := CloneSession(*s)
	ns.channelId = channelId
	return ns
}

func (s *Model) WorldId() byte {
	return s.worldId
}

func (s *Model) ChannelId() byte {
	return s.channelId
}

func (s *Model) updateLastRequest() Model {
	ns := CloneSession(*s)
	ns.lastPacket = time.Now()
	return ns
}

func (s *Model) LastRequest() time.Time {
	return s.lastPacket
}

func (s *Model) Disconnect() {
	_ = s.con.Close()
}
