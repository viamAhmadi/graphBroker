package conn

import (
	"fmt"
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
func OpenDeal(destination string) (*goczmq.Sock, error) {
	dealer, err := goczmq.NewDealer(destination)
	if err != nil {
		return nil, err
	}
	return dealer, nil
}

// SendPacketSend sends send packet
func (b *Broker) SendPacketSend(c *Connection) error {
	err := b.sock.SendFrame(c.From, goczmq.FlagMore)
	if err != nil {
		return err
	}
	return b.sock.SendFrame(SerializeSend(c.Destination, c.Sign), goczmq.FlagNone)
}

func (b *Broker) SendPacketFactor(c *Connection) error {
	err := b.sock.SendFrame(c.From, goczmq.FlagMore)
	if err != nil {
		return err
	}
	f := SerializeFactor(c.Destination, c.Sign, c.Successful, c.MissingMessages)
	fmt.Println(string(*f))
	return b.sock.SendFrame(*f, goczmq.FlagNone)
}
