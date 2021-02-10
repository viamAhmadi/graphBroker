package broker

import (
	"github.com/viamAhmadi/graphBroker/pkg/conn"
	"github.com/zeromq/goczmq"
)

type Broker struct {
	sock *goczmq.Sock
}

func New(sock *goczmq.Sock) *Broker {
	b := Broker{sock: sock}
	return &b
}

// SendPacketSend sends send packet
func (b *Broker) SendPacketSend(c *conn.Connection) error {
	err := b.sock.SendFrame(c.From, goczmq.FlagMore)
	if err != nil {
		return err
	}

	return b.sock.SendFrame(conn.SerializeSend(c.Destination, c.Sign), goczmq.FlagNone)
}
