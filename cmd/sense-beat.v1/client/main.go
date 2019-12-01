package main

import (
	"flag"
	"fmt"

	"github.com/sense-beat/cmd/sense-beat.v1/client/option"
	"github.com/sense-beat/pkg/UDPBeat"
)

func main() {
	clientOption := option.NewClientOption()
	clientOption.AddFlags(flag.CommandLine)
	flag.Parse()

	sc, err := UDPBeat.NewSockerClient(clientOption.ServerAddr, clientOption.Data, clientOption.CycleTime)
	if err != nil {
		fmt.Println(err)
	}

	sc.Serv()

}
