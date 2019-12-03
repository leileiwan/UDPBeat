package e2e

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"

	"github.com/sense-beat/pkg/UDPBeat"
	"github.com/sense-beat/pkg/watch"
)

var (
	Watcher *watch.Watcher
	SS      *UDPBeat.SocketService
	SC      *UDPBeat.SocketClient
)

func serverHttp(addr string) {
	http.HandleFunc("/getTargetStatus", func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")
		if ip == "" {
			w.Write(nil)
		}
		w.Write(Watcher.GetTargetStatus(ip))

	})
	http.HandleFunc("/getAllStatus", func(w http.ResponseWriter, r *http.Request) {
		w.Write(Watcher.GetStatusALL())
	})
	http.ListenAndServe(addr, nil)
}

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UDP Beat Test Suite")
}

var _ = BeforeSuite(func() {
	go initServer()
	go initClient()
})

var _ = AfterSuite(func() {
	By("End Suite Test")
})

func initServer() {
	Watcher = watch.NewWatcher(time.Second*1, 1, 2, 3)
	go serverHttp(":6992")
	SS, err := UDPBeat.NewSocketService("127.0.0.1:7992")
	if err != nil {
		Fail("New server error...")
		return
	}
	SS.RegConnectHandler(Watcher.Watch)
	SS.Serv()
}

func initClient() {
	SC, _ = UDPBeat.NewSockerClient("127.0.0.1:7992", "iamok", time.Second*1)

	SC.Serv()
}
