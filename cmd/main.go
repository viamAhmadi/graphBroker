package main

import (
	"log"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := application{errorLog: errorLog, infoLog: infoLog}

	app.infoLog.Println("starting broker...")

	if err := app.startBroker("tcp://127.0.0.1:5555"); err != nil {
		panic(err)
	}
}
