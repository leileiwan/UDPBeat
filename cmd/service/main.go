package main

import (
	"fmt"

	"github.com/sense-beat/pkg/UDPBeat"
)

func main() {
	host := "127.0.0.1:7788"

	ss, err := UDPBeat.NewSocketService(host)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("service running on " + host)

	ss.Serv()
}
