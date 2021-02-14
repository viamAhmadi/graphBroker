package storage

import (
	"encoding/gob"
	"fmt"
	"github.com/viamAhmadi/gReceiver2/pkg/conn"
	sConn "github.com/viamAhmadi/gSender/pkg/conn"
	"io/ioutil"
	"os"
)

func init() {
	if err := os.MkdirAll("var/connections/sent", 0770); err != nil {
		fmt.Println(err)
	}
	if err := os.MkdirAll("var/connections/received", 0770); err != nil {
		fmt.Println(err)
	}
	if err := os.MkdirAll("var/connections/failed", 0770); err != nil {
		fmt.Println(err)
	}

	if err := os.MkdirAll("var/messages/sent", 0770); err != nil {
		fmt.Println(err)
	}
	if err := os.MkdirAll("var/messages/received", 0770); err != nil {
		fmt.Println(err)
	}
	if err := os.MkdirAll("var/messages/failed", 0770); err != nil {
		fmt.Println(err)
	}
}

var YES = 1
var NO = 0

type Storage struct {
	saveSentMessages         int
	saveReceivedMessages     int
	saveFailedMessages       int
	saveSentConnection       int
	saveReceivedConnection   int
	saveFailedSentConnection int
}

func New(saveSentMessages, saveReceivedMessages, saveFailedMessages, saveSentConnection, saveReceivedConnection, saveFailedSentConnection int) *Storage {
	return &Storage{saveSentMessages, saveReceivedMessages, saveFailedMessages, saveSentConnection, saveReceivedConnection, saveFailedSentConnection}
}

//func (s *Storage) SaveSentMessages(connId string, ms *conn.Messages) error {
//	if s.saveSentMessages == NO {
//		return nil
//	}
//	f, err := os.Create(fmt.Sprintf("var/messages/sent/%s.dat", connId))
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//	return gob.NewEncoder(f).Encode(*ms)
//}
func (s *Storage) SaveReceivedMessages(connId string, ms *conn.Messages) error {
	if s.saveReceivedMessages == NO {
		return nil
	}
	f, err := os.Create(fmt.Sprintf("var/messages/recevied/%s.dat", connId))
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(*ms)
}
func (s *Storage) SaveFailedMessages(connId string, ms *conn.Messages) error {
	if s.saveFailedMessages == NO {
		return nil
	}
	f, err := os.Create(fmt.Sprintf("var/messages/failed/%s.dat", connId))
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(*ms)
}

//func (s *Storage) SaveSentConnection(c *sConn.SendConn) error {
//	if s.saveSentConnection == NO {
//		return nil
//	}
//	f, err := os.Create(fmt.Sprintf("var/connections/sent/%s.dat", c.Id))
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//	return gob.NewEncoder(f).Encode(connSend{
//		Id:         c.Id,
//		Count:      c.Count,
//		IsOpen:     c.IsOpen,
//		Successful: c.Successful,
//	})
//}
func (s *Storage) SaveReceivedConnection(c *conn.ReceiveConn) error {
	if s.saveReceivedConnection == NO {
		return nil
	}
	f, err := os.Create(fmt.Sprintf("var/connections/received/%s.dat", c.Id))
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(connSend{
		Id:         c.Id,
		Count:      c.Count,
		IsOpen:     c.IsOpen,
		Successful: c.Successful,
	})
}
func (s *Storage) SaveFailedSentConnection(c *sConn.SendConn) error {
	if s.saveFailedSentConnection == NO {
		return nil
	}
	f, err := os.Create(fmt.Sprintf("var/connections/failed/%s.dat", c.Id))
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(connSend{
		Id:         c.Id,
		Count:      c.Count,
		IsOpen:     c.IsOpen,
		Successful: c.Successful,
	})
}

func (s *Storage) ReadFailedSentMessages(connId string) (*conn.Messages, error) {
	f, err := os.Open(fmt.Sprintf("var/messages/failed/%s.dat", connId))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var ms conn.Messages
	return &ms, gob.NewDecoder(f).Decode(&ms)
}

func (s *Storage) ReadFailedSentConnections() (*sConn.SendConns, error) {
	files, err := ioutil.ReadDir("./var/connections/failed")
	if err != nil {
		return nil, err
	}

	var conns = sConn.SendConns{}
	for _, file := range files {
		f, err := os.Open(fmt.Sprintf("var/connections/failed/%s", file.Name()))
		if err != nil {
			return nil, err
		}
		var c sConn.SendConn
		if err := gob.NewDecoder(f).Decode(&c); err != nil {
			return nil, err
		}
		if err := conns.Add(&c); err != nil {
			return nil, err
		}
		f.Close()
	}
	return &conns, nil
}
