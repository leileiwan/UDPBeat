package main

import (
	"fmt"
	"time"

	"github.com/9b9387/zero/pkg/tcpBeat"
)

func HandleMessage(s *tcpBeat.Session, msg *tcpBeat.Message) {
	fmt.Println("receive msgID:", msg)
	fmt.Println("receive data:", string(msg.GetData()))
}

func HandleDisconnect(s *tcpBeat.Session, err error) {
	fmt.Println(s.GetConn().GetName() + " lost.")
}

func HandleConnect(s *tcpBeat.Session) {
	fmt.Println(s.GetConn().GetName() + " connected.")
}

func main() {
	host := "127.0.0.1:7788"

	ss, err := tcpBeat.NewSocketService(host)
	if err != nil {
		fmt.Println(err)
		return
	}

	// set Heartbeat
	ss.SetHeartBeat(5*time.Second, 30*time.Second)

	// net event
	ss.RegMessageHandler(HandleMessage)
	ss.RegConnectHandler(HandleConnect)
	ss.RegDisconnectHandler(HandleDisconnect)

	// timer := time.NewTimer(time.Second * 1)
	// go func() {
	// 	<-timer.C
	// 	ss.Stop("stop service")
	// 	fmt.Println("service stoped")
	// }()

	fmt.Println("service running on " + host)

	ss.Serv()
}
