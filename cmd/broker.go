package main

import (
	"github.com/viamAhmadi/graphBroker/pkg/conn"
	"github.com/zeromq/goczmq"
	"time"
)

func (a *application) startBroker(endpoints string) error {
	router, err := goczmq.NewRouter(endpoints)
	if err != nil {
		return err
	}

	a.broker = conn.New(router)

	for {
		msg, err := router.RecvMessage()
		if err != nil {
			a.errorLog.Println(err)
			continue
		}
		go a.router(&msg)
	}
}

func (a *application) router(rc *[][]byte) {
	valStr := string((*rc)[1][0])
	from := (*rc)[0]

	if valStr == "c" {
		go a.connectionHandler(from, &(*rc)[1])
		//a.infoLog.Println("new connection")
	} else if valStr == "m" {
		time.Sleep(3 * time.Second) // todo
		go a.messageHandler(from, &(*rc)[1])
		//a.infoLog.Println("new message")
	} else if valStr == "d" {
		go a.doneHandler(from, &(*rc)[1])
		//a.infoLog.Println("new message")
	} else {
		a.infoLog.Printf("there is unknown type, value: %v\n", valStr)
	}
}
