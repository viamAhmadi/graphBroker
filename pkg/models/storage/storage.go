package storage

import "github.com/viamAhmadi/graphBroker/pkg/conn"

type Storage struct {
}

func New(path string) Storage {
	return Storage{}
}

func (s *Storage) AddConn(conn *conn.Connection) error {
	return nil
}

func (s *Storage) AddMessage(msg *conn.Message) error {
	return nil
}
