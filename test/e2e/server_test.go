package e2e

import (
	"encoding/json"

	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sense-beat/pkg/UDPBeat"
	"github.com/sense-beat/pkg/watch"
)

type JsonMessage struct {
	IP     string
	Status string
}

var (
	clientIP  string
	statusMap map[string]string
	eachMsg   JsonMessage
	Alive     string
	NotFound  string
	Dead      string
	Watcher   *watch.Watcher
	SS        *UDPBeat.SocketService
	SC        *UDPBeat.SocketClient
)

var _ = Describe("Server Test", func() {
	BeforeEach(func() {
		Alive = "alive"
		NotFound = "notfound"
		Dead = "dead"
	})

	initHttpServer()

	Context("Get the common status", func() {
		JustBeforeEach(func() {
			statusMap = make(map[string]string, 10)
			eachMsg = JsonMessage{}
			clientIP, _ = UDPBeat.GetInternal()
			initServer()
			initClient()
		})
		JustAfterEach(func() {
			SS.Close()
			SC.Close()
			SS = nil
			SC = nil
			Watcher.Clean()

		})
		It("The Status should be alive", func() {
			time.Sleep(time.Second * 3)
			jsonRes := httpGet("http://127.0.0.1:6992/getTargetStatus?ip=" + clientIP)
			json.Unmarshal(jsonRes, &eachMsg)
			statusMap[eachMsg.IP] = eachMsg.Status
			Expect(statusMap[clientIP]).To(Equal(Alive))
		})

		It("The Status should be notfound", func() {
			time.Sleep(time.Second * 3)
			jsonRes := httpGet("http://127.0.0.1:6992/getTargetStatus?ip=127.0.0.2")
			json.Unmarshal(jsonRes, &eachMsg)
			statusMap[eachMsg.IP] = eachMsg.Status
			Expect(statusMap["127.0.0.2"]).To(Equal(NotFound))
		})
		It("The status should be dead", func() {
			SC.SetRecycleTime(time.Second * 6)
			time.Sleep(time.Second * 3)
			jsonRes := httpGet("http://127.0.0.1:6992/getTargetStatus?ip=" + clientIP)
			json.Unmarshal(jsonRes, &eachMsg)
			statusMap[eachMsg.IP] = eachMsg.Status
			Expect(statusMap[clientIP]).To(Equal(Dead))

		})

		It("The host should be delete", func() {
			SC.SetRecycleTime(time.Second * 15)
			time.Sleep(time.Second * 10)
			jsonRes := httpGet("http://127.0.0.1:6992/getTargetStatus?ip=" + clientIP)
			json.Unmarshal(jsonRes, &eachMsg)
			statusMap[eachMsg.IP] = eachMsg.Status
			Expect(statusMap[clientIP]).To(Equal(NotFound))
		})

	})

})

func httpGet(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		Fail(err.Error())

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Fail(err.Error())
	}
	return body
}

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

func initServer() {
	time.Sleep(time.Second * 2)
	SS, _ = UDPBeat.NewSocketService("127.0.0.1:7992")
	SS.RegConnectHandler(Watcher.Watch)
	go SS.Serv()

}

func initClient() {
	time.Sleep(time.Second * 2)
	SC, _ = UDPBeat.NewSockerClient("127.0.0.1:7992", "iamok", time.Second*1)
	go SC.Serv()

}

func initHttpServer() {
	Watcher = watch.NewWatcher(time.Second*1, 1, 2, 3)
	go serverHttp(":6992")

}
