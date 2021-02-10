package main

import (
	"fmt"
	"github.com/viamAhmadi/graphBroker/pkg/conn"
	"time"
)

func (a *application) newConnectionHandler(from []byte, rc *[]byte) {
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
			c.CloseConnection()
		}
	}()

	for {
		select {
		case m := <-c.ReceiveMsgCh:
			fmt.Println(m)
		case f := <-c.ReceiveFactor:
			fmt.Println(f)
		case d := <-c.ReceiveDoneCh:
			fmt.Println(d)
		case <-time.Tick(5 * time.Second):
			c.CloseConnection()
		case <-c.CloseConnCh:
			close(c.ReceiveMsgCh)
			close(c.ReceiveDoneCh)
			close(c.ReceiveFactor)
			close(c.CloseConnCh)
			break
		}
	}
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

	// to send to the destination
	c.ReceiveMsgCh <- msg

	if err := a.storage.AddMessage(msg); err != nil {
		a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
		return
	}
}
