package main

import (
	"flag"
	gConn "github.com/viamAhmadi/gReceiver2/pkg/conn"
	"github.com/viamAhmadi/graphBroker/pkg/storage"
	"log"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type application struct {
	receiver      *gConn.Receiver
	errorLog      *log.Logger
	infoLog       *log.Logger
	ReceivedConns gConn.ReceivedConns
	storage       *storage.Storage
}

func main() {
	addr := flag.String("addr", "tcp://127.0.0.1:5555", "Broker address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		ReceivedConns: gConn.ReceivedConns{},
		storage: storage.New(
			storage.NO, storage.NO, storage.YES, storage.NO, storage.NO, storage.YES,
		), // todo
	}

	app.infoLog.Println("starting broker...")

	if err := app.startReceiving(*addr); err != nil {
		panic(err)
	}
}
