package conn

type Connection struct {
	Type                 string // Type of request, c=connection or m=message ..
	From                 []byte
	Destination          string
	Sign                 string
	Count                int
	FirstMsgId, EndMsgId int
	ReceiveMsgCh         chan *Message // ReceiveMsgCh message packet
	ReceiveDoneCh        chan *Done    // ReceiveDoneCh done packet
	ReceiveFactor        chan *Factor  // ReceiveFactor factor packet
	CloseConnCh          chan struct{}
}

func (c *Connection) GetId() string {
	return c.Destination + c.Sign
}

func (c *Connection) CloseConnection() {
	c.CloseConnCh <- struct{}{}
}
