package main

import (
	"fmt"
	"github.com/viamAhmadi/graphBroker/pkg/conn"
	"time"
)

func (a *application) newReceiverConnectionHandler(from []byte, rc *[]byte) {
	c, err := conn.ConvertToConnection(from, *rc)
	if err != nil {
		a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
		return
	}
	// memory tmp
	if err := a.conns.Add(c); err != nil {
		a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
		return
	}

	if err := a.storage.AddConn(c); err != nil {
		a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
		return
	}
	go func() {
		if err := a.broker.SendPacketSend(c); err != nil {
			a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
		}
	}()
	for {
		select {
		case m := <-c.ReceiveMsgCh:
			fmt.Println(m)
		case d := <-c.ReceiveDoneCh:
			fmt.Println(d)
			if err := a.broker.SendPacketFactor(); err != nil {
				a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
			}
		case <-time.After(5 * time.Second):
			c.CloseConnection()
		case <-c.CloseConnCh:
			close(c.ReceiveMsgCh)
			close(c.ReceiveDoneCh)
			close(c.ReceiveFactor)
			return
		}
	}
}

// todo
func (a *application) newSenderConnectionHandler() {

}

func (a *application) newMessageHandler(from []byte, rc *[]byte) {
	msg, err := conn.ConvertToMessage(rc)
	if err != nil {
		a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
		return
	}
	c := a.conns.Get(msg.GetConnId())
	if c == nil {
		a.sendPacketError(conn.Error{Msg: "connection not found", Destination: from})
		return
	}

	c.ReceiveMsgCh <- msg



	if err := a.storage.AddMessage(msg); err != nil {
		a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
		return
	}
}
