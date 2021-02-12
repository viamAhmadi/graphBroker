package main

import (
	"fmt"
	"github.com/viamAhmadi/graphBroker/pkg/conn"
	"runtime/debug"
)

func (a *application) sendPacketError(p conn.Error) {
	a.errorLog.Println(fmt.Sprintf("%s\n%s", p.Msg, debug.Stack()))
}
