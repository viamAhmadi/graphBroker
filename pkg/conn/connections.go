package conn

import (
	"errors"
)

var (
	Open  = uint8(1)
	Close = uint8(0)
)

const ok = "1"
const YES = byte(1)
const NO = byte(0)

var ErrConvertToModel = errors.New("convert error")
var ErrConnExists = errors.New("connection exists")
var ErrDealer = errors.New("dealer was nil")

type Connections map[string]*Connection

func (c *Connections) Add(conn *Connection) error {
	if cFound := c.Get(conn.GetId()); cFound != nil {
		return ErrConnExists
	}
	(*c)[conn.GetId()] = conn
	return nil
}

func (c *Connections) Get(connId string) *Connection {
	return (*c)[connId]
}
