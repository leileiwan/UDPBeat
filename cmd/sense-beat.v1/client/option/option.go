package option

import (
	"flag"
	"time"
)

type ClientOption struct {
	ServerAddr string
	Data       string
	CycleTime  time.Duration
}

func NewClientOption() *ClientOption {
	return &ClientOption{}
}

func (c *ClientOption) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.ServerAddr, "ServerAddr", "127.0.0.1:7991", "The master url with ip and port")
	fs.StringVar(&c.Data, "Data", "I am ok", "The message that client sent to message")
	fs.DurationVar(&c.CycleTime, "CycleTime", 5*time.Second, "The cycle time that client sent the message to master")
}
