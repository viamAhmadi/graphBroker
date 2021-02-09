package main

import (
	"github.com/zeromq/goczmq"
)

func (a *application) startBroker(endpoints string) error {
	router, err := goczmq.NewRouter(endpoints)
	if err != nil {
		return err
	}
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
	valStr := string((*rc)[1])
	//from := (*rc)[0]

	if valStr == "c" {
		go a.newConnectionHandler(&(*rc)[1])
		a.infoLog.Println("new connection")
	} else if valStr == "m" {
		go a.newMessageHandler(&(*rc)[1])
		a.infoLog.Println("new message")
	} else {
		a.infoLog.Printf("there is unkown type, vlaeue: %v\n", valStr)
	}
}
