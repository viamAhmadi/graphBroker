package main

import (
	"fmt"
	"github.com/viamAhmadi/gReceiver2/pkg/conn"
	"runtime/debug"
)

func (a *application) sendPacketError(r conn.Error) {
	a.errorLog.Println(fmt.Sprintf("%s\n%s", r.Msg, debug.Stack()))
}
