package main

import "github.com/viamAhmadi/graphBroker/pkg/conn"

func (a *application) sendPacketError(p conn.Error) {
	a.errorLog.Println(p)
	// TODO send error packet, dont use concurrent
}
