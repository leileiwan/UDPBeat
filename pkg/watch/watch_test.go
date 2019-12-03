package watch

import (
	"github.com/sense-beat/pkg/UDPBeat"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var watch *Watcher

func TestMain(m *testing.M) {
	Init()
	m.Run()
}

func Init() {
	watch = NewWatcher(1, 1, 2, 3)
	host := Host{
		IP:    "127.0.0.1",
		Time:  time.Now(),
		HP:    3,
		Alive: true,
	}
	watch.hosts["127.0.0.1"] = &host
}

func TestUpdateState(t *testing.T) {
	watch.hosts["127.0.0.1"].HP = 0
	watch.updateState("127.0.0.1")
	assert.False(t, watch.hosts["127.0.0.1"].Alive)

	watch.hosts["127.0.0.1"].HP = 3
	watch.updateState("127.0.0.1")
	assert.True(t, watch.hosts["127.0.0.1"].Alive)
}

func TestHurt(t *testing.T) {
	watch.hurt("127.0.0.1")
	assert.True(t, watch.hosts["127.0.0.1"].Alive)
	assert.Equal(t, 2, watch.hosts["127.0.0.1"].HP)

	watch.hosts["127.0.0.1"].HP = 1
	watch.hurt("127.0.0.1")
	assert.False(t, watch.hosts["127.0.0.1"].Alive)
	assert.Equal(t, 0, watch.hosts["127.0.0.1"].HP)

	watch.hosts["127.0.0.1"].HP = 0
	watch.hurt("127.0.0.1")
	assert.False(t, watch.hosts["127.0.0.1"].Alive)
	assert.Equal(t, 0, watch.hosts["127.0.0.1"].HP)
}

func TestFix(t *testing.T) {
	watch.hurt("127.0.0.1")
	watch.hosts["127.0.0.1"].HP = 0
	watch.updateState("127.0.0.1")

	msg := UDPBeat.NewMessage("127.0.0.1", "iamok")
	watch.hosts["127.0.0.1"].HP = -1
	watch.fix(*msg)
	assert.Equal(t, watch.hosts["127.0.0.1"].HP, watch.levelInitHP)
	assert.False(t, watch.hosts["127.0.0.1"].Alive)

	watch.fix(*msg)
	assert.Equal(t, watch.hosts["127.0.0.1"].HP, watch.levelAliveHP)
	assert.True(t, watch.hosts["127.0.0.1"].Alive)

	watch.hosts["127.0.0.1"].HP = 3
	watch.fix(*msg)
	assert.Equal(t, watch.hosts["127.0.0.1"].HP, watch.levelFullHP)
	assert.True(t, watch.hosts["127.0.0.1"].Alive)

}
