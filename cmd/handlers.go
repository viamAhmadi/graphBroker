package main

import (
	"fmt"
	conn2 "github.com/viamAhmadi/graphBroker/pkg/conn"
)

func (a *application) newConnectionHandler(rc *[]byte) {
	conn, err := conn2.ConvertToConnection(*rc)
	if err != nil {
		a.errorLog.Println(err)
		return
	}
	fmt.Println("type ", conn.Type)
	fmt.Println("destination ", conn.Destination)
	fmt.Println("count", conn.Count)
	fmt.Println("sign", conn.Sign)
	fmt.Println("firstMsgId ", conn.FirstMsgId)
	fmt.Println("endMsgId ", conn.EndMsgId)
}

func (a *application) newMessageHandler(rc *[]byte) {

}
