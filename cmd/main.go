package main

import (
	"github.com/viamAhmadi/graphBroker/pkg/broker"
	"github.com/viamAhmadi/graphBroker/pkg/conn"
	"github.com/viamAhmadi/graphBroker/pkg/models/storage"
	"log"
	"os"
)

type application struct {
	broker   *broker.Broker
	errorLog *log.Logger
	infoLog  *log.Logger
	conns    conn.Connections
	storage  storage.Storage
}

func main() {
	//c := conn.SerializeConnection("192.168.1.1:5", "as", 1234, 1234, 4321)
	//fmt.Println(len(c))
	//fmt.Println(conn.ConvertToConnection(c))
	//
	//m := conn.SerializeMessage(1, "si", "hello")
	//fmt.Println(len(*m))
	//fmt.Println(conn.ConvertToMessage(m))
	//
	//d := conn.SerializeDone("192.168.2.3:123", "bi", 14)
	//fmt.Println(len(d))
	//fmt.Println(conn.ConvertToDone(d))
	//
	//s := conn.SerializeSend("192.168.2.5:43215", "ds")
	//fmt.Println(len(s))
	//fmt.Println(conn.ConvertToSend(s))
	//
	//list := []string{"5", "321", "789"}
	//f := conn.SerializeFactor("127.0.0.1:65432", "iq", "n", &list)
	//fmt.Println(len(*f))
	//fmt.Println(conn.ConvertToFactor(f))

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := application{
		errorLog: errorLog,
		infoLog:  infoLog,
		conns:    conn.Connections{},
		storage:  storage.New(""), // todo path of storage
	}

	app.infoLog.Println("starting broker...")

	if err := app.startBroker("tcp://127.0.0.1:5555"); err != nil {
		panic(err)
	}
}
