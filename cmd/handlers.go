package main

import (
	"fmt"
	"github.com/viamAhmadi/graphBroker/pkg/conn"
	"time"
)

func (a *application) connectionHandler(from []byte, rc *[]byte) {
	c, err := conn.ConvertToConnection(from, *rc)
	if err != nil {
		a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
		return
	}
	// memory tmp
	if err := a.conns.Add(c); err != nil {
		a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
		return
	} else {
		if err := a.storage.AddConn(c); err != nil {
			a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
			return
		}
	}

	go func() {
		// forward to the destination
		if c.Forward == conn.YES {
			fmt.Println("Broker: opening connection to A")
			dealer, err := conn.OpenDeal(c.Destination)
			if err != nil {
				a.errorLog.Println(err)
				c.Dealer = nil
			} else {
				c.Dealer = dealer
				_ = conn.SendPacketConnection(c, conn.NO)
				go a.forwardResultHandler(c)
			}
		}
		// tell the sender Im ready for receive
		fmt.Println("Are you ready? Broker and A?")
		if err := a.broker.SendPacketSend(c); err != nil {
			a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
		}
	}()
	for {
		select {
		case m := <-c.ReceiveMsgCh:
			fmt.Println("Broker: forwarding one message to A")
			m.Forward = conn.NO
			m.Sign = c.SendSign
			if err := conn.SendFrame(c, m); err != nil {
				a.errorLog.Println(err)
			}
			// check last message id and send done packet
			if c.EndMsgId == m.Id {
				fmt.Println("Broker: sending done packet to A")
				if err := conn.SendPacketDone(c, conn.NO); err != nil {
					a.errorLog.Println(err)
				}
			}
		case _ = <-c.ReceiveDoneCh:
			if c.Count == c.CountOfRcMessages() {
				c.Successful = conn.YES
			} else {
				if count := c.CalculateMissingMessages(); count != 0 {
					c.Successful = conn.NO
				}
			}
			if err := a.broker.SendPacketFactor(c); err != nil {
				a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
			}
		case <-time.After(5 * time.Second):
			c.CloseConnection()
		case <-c.CloseConnCh:
			close(c.ReceiveMsgCh)
			close(c.ReceiveDoneCh)
			close(c.ReceiveFactorCh)
			return
		}
	}
}

func (a *application) forwardResultHandler(c *conn.Connection) {
	go func() {
		for {
			select {
			case s := <-c.ReceiveSendCh:
				// if there is some errors -> dealer=nil
				// send messages , the server is ready
				//fmt.Println(s)
				// start sending ...
				fmt.Println("C[s]: I'm ready to receive  ", s)
			case f := <-c.ReceiveFactorCh:
				fmt.Println("f")
				fmt.Println(f)
				return
			case <-time.After(6 * time.Second):
				return
			}
		}
	}()
	if err := c.DealerStartReceiving(); err != nil {
		a.errorLog.Println(err)
	}
}

func (a *application) messageHandler(from []byte, rc *[]byte) {
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

	if c.Forward == conn.YES {
		if err := c.AddMsg(msg); err != nil {
			a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
		}
		c.ReceiveMsgCh <- msg
	}

	if err := a.storage.AddMessage(msg); err != nil {
		a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
		return
	}
}

func (a *application) doneHandler(from []byte, rc *[]byte) {
	dPacket, err := conn.ConvertToDone(*rc)
	if err != nil {
		a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
		return
	}

	c := a.conns.Get(dPacket.GetConnId())
	if c == nil {
		a.sendPacketError(conn.Error{Msg: "connection not found", Destination: from})
		return
	}

	fmt.Println(dPacket)
	c.ReceiveDoneCh <- &dPacket
}
