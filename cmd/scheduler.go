package main

import (
	"github.com/jasonlvhit/gocron"
	"github.com/viamAhmadi/gReceiver2/pkg/util"
)

func (a *application) StartScheduler() error {
	s1 := gocron.NewScheduler()
	if err := s1.Every(10).Minutes().Do(removeClosedConnections, a); err != nil {
		a.errorLog.Println(err)
	}
	if err := s1.Every(10).Minutes().Do(tryToForwardFailedMessages, a); err != nil {
		a.errorLog.Println(err)
	}

	<-s1.Start()
	a.errorLog.Println("scheduler stopped")
	return nil
}

func removeClosedConnections(a *application) {
	for _, conn := range a.ReceivedConns {
		if conn.IsOpen == 2 { // 2 = undefined
			continue
		}
		a.ReceivedConns.Remove(conn.Id)
	}
}

func tryToForwardFailedMessages(a *application) {
	conns, err := a.storage.ReadFailedSentConnections()
	if err != nil {
		a.errorLog.Println(err)
		return
	}
	go func() {
		for _, conn := range *conns {
			messages, err := a.storage.ReadFailedSentMessages(conn.Id)
			if err != nil {
				a.errorLog.Println(err)
				continue
			}
			conn.Count = len(*messages)
			conn.Id = util.RandomString(20)
			conn.Messages = messages
			a.forward(conn)
		}
	}()
}
