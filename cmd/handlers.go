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
	fmt.Println("connId: ", c.GetId())
	// memory tmp
	if err := a.conns.Add(c); err != nil {
		a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
	}
	if err := a.storage.AddConn(c); err != nil {
		a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
		return
	}

	go func() {
		// forward to the destination
		if c.Forward == conn.YES {
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
		if err := a.broker.SendPacketSend(c); err != nil {
			a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
		}
	}()
	r := false
	for {
		select {
		case m := <-c.ReceiveMsgCh:
			if !r {
				r = true
				time.Sleep(1 * time.Second)
			}
			m.Forward = conn.NO
			m.Sign = c.SendSign
			if err := conn.SendFrame(c, m); err != nil {
				a.errorLog.Println(err)
			}
			if c.Count == m.Id {
				if err := conn.SendPacketDone(c, conn.NO); err != nil {
					a.errorLog.Println(err)
				}
			}
		case _ = <-c.ReceiveDoneCh:
			fmt.Println("receive done packet")
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
				fmt.Println("ReceiveSendCh  ", s)
			case f := <-c.ReceiveFactorCh:
				if f == nil {
					fmt.Println("factor was nil")
					return
				}
				l := 0
				if f.List != nil {
					l = len(f.List)
				}
				fmt.Printf("\nReceived Factor\nresult:\t\tsend: %d\tmissed: %d \ndetination: %s\n\n", c.Count, l, f.Destination)
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
		c.ReceiveMsgCh <- msg
	}

	if err := c.AddMsg(msg); err != nil {
		a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
	}

	if c.Count == c.Counter {
		c.Successful = conn.YES
		if err := a.broker.SendPacketFactor(c); err != nil {
			a.sendPacketError(conn.Error{Msg: err.Error(), Destination: from})
		}
		c.CloseConnection()
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

	if c.IsClosed != 0 {
		c.ReceiveDoneCh <- &dPacket
	}
}
