package main

import (
	"fmt"
	gConn "github.com/viamAhmadi/gReceiver2/pkg/conn"
	"github.com/viamAhmadi/gReceiver2/pkg/util"
	sConn "github.com/viamAhmadi/gSender/pkg/conn"
	"strconv"
	"time"
)

func (a *application) forward(cRec *gConn.ReceiveConn) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	if cRec.MissedCount > 500 {
		fmt.Println("\n* missed ", cRec.MissedCount)
		return
	}
	c := sConn.NewSendConn(cRec.Destination, cRec.Id, cRec.Count)
	if err := c.OpenConnection(cRec.Destination); err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case _ = <-c.SendCh:
				fmt.Printf("\nFORWARD %d", c.Count)
				if err := sendTmpMessages(cRec, c); err != nil {
					fmt.Println("error in forwarding")
				}
			case f := <-c.FactorCh:
				c.Factor = f
				missed := len(*f.List)
				result := "NO"
				if f.Successful == gConn.YES {
					result = "YES"
				} else {
					if missed >= 500 {
						go func() {
							messages := gConn.Messages{}
							for _, msgId := range *f.List {
								msgIdInt, err := strconv.Atoi(msgId)
								if err != nil {
									continue
								}
								_ = messages.Add(cRec.Messages.Get(msgIdInt))
							}
							if err := a.storage.SaveFailedMessages(c.Id, &messages); err != nil {
								a.errorLog.Println(err)
							}
							if err := a.storage.SaveFailedSentConnection(cRec); err != nil {
								a.errorLog.Println(err)
							}
						}()
					}
				}
				fmt.Printf("\nFACTOR\t\tsuccessful: %s\tlost: %d\n", result, missed)
				//c.Close()
			case <-time.After(util.CalculateTimeout(4, c.Count)): // 4
				if c.Factor == nil {
					fmt.Println("\nFAILED")
					go func() {
						if err := a.storage.SaveFailedMessages(c.Id, cRec.Messages); err != nil {
							a.errorLog.Println(err)
						}
						if err := a.storage.SaveFailedSentConnection(cRec); err != nil {
							a.errorLog.Println(err)
						}
					}()
				}
				c.Close()
				return
			}
		}
	}()

	if err := c.Receiving(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("exit")
}

func sendTmpMessages(cReceive *gConn.ReceiveConn, cSend *sConn.SendConn) error {
	counter := 1
	//for _, val := range *cReceive.Messages {
	//	if counter == 30000 {
	//		time.Sleep(1 * time.Second)
	//		counter = 1
	//	}
	//	counter += 1
	//	if err := cSend.Send(val); err != nil {
	//		fmt.Println(err)
	//	}
	//}
	for i := 1; i <= cReceive.Count; i++ {
		msg, ok := (*cReceive.Messages)[i]
		if !ok {
			continue
		}
		if counter == 30000 {
			time.Sleep(1 * time.Second)
			counter = 1
		}
		counter += 1
		if err := cSend.Send(msg); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
