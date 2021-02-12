package conn

import (
	"errors"
	"sync"
)

var (
	ErrMsgExist = errors.New("message exist")
)

type Messages map[string]*Message

var mutexM = sync.Mutex{}

func (m *Messages) Add(msg *Message) error {
	if mFound := m.Get(msg.GetId()); mFound != nil {
		return ErrMsgExist
	}
	mutexM.Lock()
	(*m)[msg.GetId()] = msg
	mutexM.Unlock()
	return nil
}

func (m *Messages) Get(msgId string) *Message {
	mutexM.Lock()
	msg := (*m)[msgId]
	mutexM.Unlock()
	return msg
}

func (m *Messages) Count() int {
	return len(*m)
}
