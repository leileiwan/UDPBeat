package UDPBeat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAndDecode(t *testing.T) {
	msg := NewMessage("127.0.0.1", "I am ok")
	encodeBytes, _ := Encode(msg)
	decodeMsg, _ := Decode(encodeBytes)
	assert.Equal(t, msg.GetIP(), decodeMsg.GetIP())
	assert.Equal(t, msg.GetData(), decodeMsg.GetData())
}

func TestHostAddrCheck(t *testing.T) {
	res := HostAddrCheck("127.0.0.1:111")
	assert.True(t, res)

	res = HostAddrCheck("1d7.0.0.1:111")
	assert.False(t, res)

	res = HostAddrCheck("127.0.0.1:1a1")
	assert.False(t, res)

	res = HostAddrCheck(":111")
	assert.False(t, res)

	res = HostAddrCheck("127.0.0.1:")
	assert.False(t, res)

	res = HostAddrCheck("127.0..1:")
	assert.False(t, res)
	res = HostAddrCheck("127.0.1:")
	assert.False(t, res)

	res = HostAddrCheck("1271.0.0.1:111")
	assert.False(t, res)

	res = HostAddrCheck("127.0.0.1:11111111111")
	assert.False(t, res)

	res = HostAddrCheck("127.0.0.1.1:111")
	assert.False(t, res)

}
