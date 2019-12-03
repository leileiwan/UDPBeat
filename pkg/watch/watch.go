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

// use udp protocal
type Watcher struct {
	sync.Mutex
	recycleDuration time.Duration
	hosts           map[string]*Host
	levelInitHP     int // 接收到心跳时，HP最低为该值
	levelAliveHP    int // 复活时的HP
	levelFullHP     int // 满血时的HP
}

func NewWatcher(cycletime time.Duration, levelInitHP, levelAliveHP, levelFullHP int) *Watcher {
	return &Watcher{
		recycleDuration: cycletime,
		hosts:           make(map[string]*Host, 100),
		levelInitHP:     levelInitHP,
		levelAliveHP:    levelAliveHP,
		levelFullHP:     levelFullHP,
	}
}

// cut down HP
func (this *Watcher) hurt(ip string) {
	h, ok := this.hosts[ip]
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
	h, ok := this.hosts[ip]
	if !ok {
		this.hosts[ip] = &Host{IP: ip, Time: time.Now(), HP: this.levelFullHP, Alive: true}
		return
	}
	h.HP += 1
	if h.HP > this.levelFullHP {
		h.HP = this.levelFullHP
	}
	if h.HP < this.levelInitHP {
		h.HP = this.levelInitHP
	}
	h.Time = time.Now()
	this.updateState(ip)
}

// judge if host is Alive from HP
func (this *Watcher) updateState(ip string) {
	host, ok := this.hosts[ip]
	if !ok {
		return
	}
	if host.HP >= this.levelAliveHP {
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
	alives := make([]string, 0, len(this.hosts))
	deads := make([]string, 0, len(this.hosts))
	for ip, host := range this.hosts {
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
	host, ok := this.hosts[ip]
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
		for _, host := range this.hosts {
			this.hurt(host.IP)
		}
		this.Unlock()
		time.Sleep(this.recycleDuration)
	}
}

//add for test
func (watcher *Watcher) SetRecycleTime(time time.Duration) {
	watcher.recycleDuration = time
}
