package conn

import (
	"fmt"
	"github.com/viamAhmadi/graphBroker/pkg/util"
	"github.com/zeromq/goczmq"
	"strconv"
)

type Connection struct {
	Type                 string // Type of request, c=connection or m=message ..
	Forward              byte   // 1 - 0
	From                 []byte
	Destination          string
	Sign                 string
	Count                int
	FirstMsgId, EndMsgId int
	ReceiveMsgCh         chan *Message // ReceiveMsgCh message packet
	ReceiveDoneCh        chan *Done    // ReceiveDoneCh done packet
	ReceiveSendCh        chan *Send    // ReceiveSendCh send packet
	ReceiveFactorCh      chan *Factor  // ReceiveFactorCh factor packet
	CloseConnCh          chan struct{}
	Successful           byte
	Messages             *Messages
	MissingMessages      *[]string
	Dealer               *goczmq.Sock
	SendSign             string
	IsClosed             int
	Counter              int
}

func ConvertToConnection(from []byte, b []byte) (*Connection, error) {
	if cap(b) < 39 {
		return nil, ErrConvertToModel
	}
	count, err := strconv.Atoi(util.RemoveAdditionalCharacters(b[27:31]))
	if err != nil {
		return nil, err
	}
	firstMsgId, err := strconv.Atoi(util.RemoveAdditionalCharacters(b[31:35]))
	if err != nil {
		return nil, err
	}
	endMsgId, err := strconv.Atoi(util.RemoveAdditionalCharacters(b[35:39]))
	if err != nil {
		return nil, err
	}
	return &Connection{
		Type:            string(b[0]),
		Forward:         b[39],
		From:            from,
		Destination:     util.RemoveAdditionalCharacters(b[1:23]),
		Sign:            util.RemoveAdditionalCharacters(b[23:27]),
		Count:           count,
		FirstMsgId:      firstMsgId,
		EndMsgId:        endMsgId,
		ReceiveMsgCh:    make(chan *Message,5000),
		ReceiveDoneCh:   make(chan *Done),
		ReceiveSendCh:   make(chan *Send),
		ReceiveFactorCh: make(chan *Factor),
		CloseConnCh:     make(chan struct{}),
		Messages:        &Messages{},
		MissingMessages: &[]string{},
		Successful:      NO,
		SendSign:        "-1",
		IsClosed:        0,
		Counter:         0,
	}, nil
}

func SerializeConnection(forward byte, destination, sign string, count, firstMsgId, endMsgId int) []byte {
	s, _ := strconv.Atoi(sign)
	b := []byte(fmt.Sprintf("c%s%s%s%s%s", util.ConvertDesToBytes(destination), util.ConvertIntToBytes(s), util.ConvertIntToBytes(count), util.ConvertIntToBytes(firstMsgId), util.ConvertIntToBytes(endMsgId)))
	return append(b, forward)
}

func (c *Connection) CountOfRcMessages() int {
	return c.Messages.Count()
}

func (c *Connection) CalculateMissingMessages() int {
	var missed int
	for i := c.FirstMsgId; i <= c.EndMsgId; i++ {
		strI := strconv.Itoa(i)
		if m := c.Messages.Get(strI); m == nil {
			*c.MissingMessages = append(*c.MissingMessages, strI)
			missed += 1
		}
	}
	return missed
}

func (c *Connection) GetId() string {
	if c.Forward == YES {
		return c.Destination + c.Sign + "s"
	}
	return c.Destination + c.Sign
}

func (c *Connection) AddMsg(m *Message) error {
	return c.Messages.Add(m)
}

func (c *Connection) CloseConnection() {
	if c.IsClosed != 1 {
		close(c.CloseConnCh)
		c.IsClosed = 1
	}
}

func SendFrame(c *Connection, msg *Message) error {
	if c.Dealer == nil {
		return ErrDealer
	}
	// TODO I should use flagMore and check last message id
	//fmt.Println(*msg)
	return c.Dealer.SendFrame(*SerializeMessage(msg.Id, msg.Forward, msg.Sign, msg.Destination, &msg.Content), goczmq.FlagNone)
}

func SendPacketConnection(c *Connection, forward byte) error {
	if c.Dealer == nil {
		return ErrDealer
	}
	sign := generateSign(c.Destination)
	cPacket := SerializeConnection(forward, c.Destination, sign, c.Count, c.FirstMsgId, c.EndMsgId)
	c.SendSign = sign
	return c.Dealer.SendFrame(cPacket, goczmq.FlagNone)
}

func SendPacketDone(c *Connection, forward byte) error {
	if c.Dealer == nil {
		return ErrDealer
	}
	return c.Dealer.SendFrame(SerializeDone(c.Destination, c.SendSign, forward, c.Count), goczmq.FlagNone)
}

func (c *Connection) DealerStartReceiving() error {
	if c.Dealer == nil {
		return ErrDealer
	}

	for {
		msg, err := c.Dealer.RecvMessage()
		if err != nil {
			return err
		}

		valStr := string((msg)[0][0])

		if valStr == "s" {
			s, err := ConvertToSend(msg[0])
			if err != nil {
				return err
			}
			c.ReceiveSendCh <- &s
		} else if valStr == "f" {
			f, err := ConvertToFactor(&(msg)[0])
			if err != nil {
				return err
			}
			c.ReceiveFactorCh <- f
			c.Dealer.Destroy()
			return nil
		}
	}
}
