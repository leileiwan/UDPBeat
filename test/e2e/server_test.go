package e2e

import (
	"encoding/json"

	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sense-beat/pkg/UDPBeat"
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
)

var _ = Describe("Server Test", func() {
	BeforeEach(func() {
		Alive = "alive"
		NotFound = "notfound"
		Dead = "dead"
	})

	Context("Get the common status", func() {
		JustBeforeEach(func() {
			statusMap = make(map[string]string, 10)
			eachMsg = JsonMessage{}
			clientIP, _ = UDPBeat.GetInternal()
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
			SC.SetRecycleTime(time.Second * 5)
			time.Sleep(time.Second * 6)
			jsonRes := httpGet("http://127.0.0.1:6992/getTargetStatus?ip=" + clientIP)
			json.Unmarshal(jsonRes, &eachMsg)
			statusMap[eachMsg.IP] = eachMsg.Status
			Expect(statusMap[clientIP]).To(Equal(Dead))

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
