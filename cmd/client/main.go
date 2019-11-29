package main

import (
	"net"

	"time"

	"github.com/9b9387/zero/pkg/tcpBeat"
)

func main() {
	go NewClientConnect()
	time.Sleep(time.Second * 2)

}
func NewClientConnect() {
	host := "127.0.0.1:18787"
	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return
	}

	msg := tcpBeat.NewMessage(1, []byte("Hello Zero!"))
	data, err := tcpBeat.Encode(msg)
	if err != nil {
		return
	}
	conn.Write(data)
}
