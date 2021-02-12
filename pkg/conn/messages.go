package conn

import "errors"

var (
	ErrMsgExist = errors.New("message exist")
)

type Messages map[string]*Message

func (m *Messages) Add(msg *Message) error {
	if mFound := m.Get(msg.GetId()); mFound != nil {
		return ErrMsgExist
	}
	(*m)[msg.GetId()] = msg
	return nil
}

func (m *Messages) Get(msgId string) *Message {
	return (*m)[msgId]
}

func (m *Messages) Count() int {
	return len(*m)
}
