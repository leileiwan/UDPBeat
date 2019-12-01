package main

import (
	"flag"
	"fmt"

	"net/http"

	"github.com/sense-beat/cmd/sense-beat.v1/server/option"
	"github.com/sense-beat/pkg/UDPBeat"
	"github.com/sense-beat/pkg/watch"
)

var watcher = watch.NewWatcher()

func serverHttp(addr string) {
	http.HandleFunc("/getTargetStatus", func(w http.ResponseWriter, r *http.Request) {

		ip := r.URL.Query().Get("ip")
		if ip == "" {
			fmt.Printf("Not found the para ip with the request %v\n", r)
		}
		w.Write(watcher.GetTargetStatus(ip))

	})
	http.HandleFunc("/getAllStatus", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		w.Write(watcher.GetStatusALL())
	})

	http.ListenAndServe(addr, nil)
}

func main() {
	serverOption := option.NewServerOption()
	serverOption.AddFlags(flag.CommandLine)
	flag.Parse()

	go serverHttp(":" + serverOption.HttpPort)

	host := "127.0.0.1:" + serverOption.ListenPort

	ss, err := UDPBeat.NewSocketService(host)
	if err != nil {
		fmt.Println(err)
		return
	}

	ss.RegConnectHandler(watcher.Watch)
	fmt.Println("service running on " + host)

	ss.Serv()

}
