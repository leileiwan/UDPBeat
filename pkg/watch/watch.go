package watch

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/sense-beat/pkg/UDPBeat"
)

type Host struct {
	IP    string
	Time  time.Time
	HP    int
	Alive bool
}

var (
	LevelInitHP    int = 3  // 接收到心跳时，HP最低为该值
	LevelAliveHP   int = 7  // 复活时的HP
	LevelFullHP    int = 10 // 满血时的HP
	DEFAULTRECYCLE     = time.Second * 2
)

// use udp protocal
type Watcher struct {
	sync.Mutex
	RecycleDuration time.Duration
	Hosts           map[string]*Host
}

func NewWatcher() *Watcher {
	return &Watcher{
		RecycleDuration: DEFAULTRECYCLE,
		Hosts:           make(map[string]*Host, 100),
	}
}

// cut down HP
func (this *Watcher) hurt(ip string) {
	h, ok := this.Hosts[ip]
	if !ok {
		return
	}
	if h.HP -= 1; h.HP < 0 {
		h.HP = 0
	}
	this.updateState(ip)
}

// recover HP
func (this *Watcher) fix(msg UDPBeat.Message) {
	this.Lock()
	defer this.Unlock()
	ip := msg.GetIP()
	h, ok := this.Hosts[ip]
	if !ok {
		this.Hosts[ip] = &Host{IP: ip, Time: time.Now(), HP: LevelFullHP, Alive: true}
		return
	}
	h.HP += 1
	if h.HP > LevelFullHP {
		h.HP = LevelFullHP
	}
	if h.HP < LevelInitHP {
		h.HP = LevelInitHP
	}
	h.Time = time.Now()
	this.updateState(ip)
}

// judge if host is Alive from HP
func (this *Watcher) updateState(ip string) {
	host, ok := this.Hosts[ip]
	if !ok {
		return
	}
	if host.HP >= LevelAliveHP {
		host.Alive = true
	}
	if host.HP == 0 {
		host.Alive = false
	}
}

// return {"Alives": [...], "Deads": [...]}
func (this *Watcher) GetStatusALL() []byte {
	this.Lock()
	defer this.Unlock()
	alives := make([]string, 0, len(this.Hosts))
	deads := make([]string, 0, len(this.Hosts))
	for ip, host := range this.Hosts {
		if host.Alive {
			alives = append(alives, ip)
		} else {
			deads = append(deads, ip)
		}
	}
	sort.Strings(alives)
	sort.Strings(deads)
	data, err := json.Marshal(struct {
		Alives []string
		Deads  []string
	}{
		alives,
		deads,
	})
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func (this *Watcher) GetTargetStatus(ip string) []byte {
	var data struct {
		IP     string
		Status string //alive,dead,notfound
	}
	data.IP = ip
	host, ok := this.Hosts[ip]
	if !ok {
		data.Status = "notfound"
	} else {
		if host.Alive {
			data.Status = "alive"
		} else {
			data.Status = "dead"
		}

	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	return dataBytes
}

// When host is dead or come to alive, chan calls.
func (this *Watcher) Watch(ch chan UDPBeat.Message) {
	go func() {
		for {
			msg := <-ch
			this.fix(msg)
		}
	}()

	go this.drain() // clean program
}

// auto decrease host HP
func (this *Watcher) drain() {
	for {
		this.Lock()
		for _, host := range this.Hosts {
			this.hurt(host.IP)
		}
		this.Unlock()
		time.Sleep(this.RecycleDuration)
	}
}
