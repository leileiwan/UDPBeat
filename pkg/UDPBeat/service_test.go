package UDPBeat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSocketService(t *testing.T) {
	_, err := NewSocketService("127.0.0.1:6991")
	assert.Equal(t, err, nil)
	_, err = NewSocketService("127.0.0.1:199")
	assert.NotEqual(t, err, nil)
	_, err = NewSocketService("127.0.0.1:19999999")
	assert.NotEqual(t, err, nil)
}
