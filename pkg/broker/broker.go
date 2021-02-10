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

// todo
func (b *Broker) OpenConnection(destination, destinationOfPackets string, count, firstMsgId, endMsgId int) (*goczmq.Sock, error) {
	cPacket := conn.SerializeConnection(destinationOfPackets, generateSign(destinationOfPackets), count, firstMsgId, endMsgId)
	dealer, err := goczmq.NewDealer(destination)
	if err != nil {
		return nil, err
	}
	return dealer, dealer.SendFrame(cPacket, goczmq.FlagNone)
}

// SendPacketSend sends send packet
func (b *Broker) SendPacketSend(c *conn.Connection) error {
	err := b.sock.SendFrame(c.From, goczmq.FlagMore)
	if err != nil {
		return err
	}

	return b.sock.SendFrame(conn.SerializeSend(c.Destination, c.Sign), goczmq.FlagNone)
}

// todo
func (b *Broker) SendPacketFactor() error {
	return nil
}
