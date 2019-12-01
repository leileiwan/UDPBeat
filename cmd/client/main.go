package main

import (
	"flag"
	"fmt"

	"github.com/sense-beat/pkg/UDPBeat"
)

var (
	serverAddr = flag.String("serverAddr", "127.0.0.1:7788", "heartbeat server address")
	data       = flag.String("data", "I am ok", "The message sent to server")
	cycleTime  = flag.Int("cycleTime", 5, "The cycle time that sent a message to server")
)

func main() {
	flag.Parse()
	sc, err := UDPBeat.NewSockerClient(*serverAddr, *data, *cycleTime)
	if err != nil {
		fmt.Println(err)
	}

	sc.Serv()

}
