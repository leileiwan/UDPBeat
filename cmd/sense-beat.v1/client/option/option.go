package option

import (
	"flag"
)

type ClientOption struct {
	ServerAddr string
	Data       string
	CycleTime  int
}

func NewClientOption() *ClientOption {
	return &ClientOption{}
}

func (c *ClientOption) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.ServerAddr, "ServerAddr", "127.0.0.1:7991", "The master url with ip and port")
	fs.StringVar(&c.Data, "Data", "I am ok", "The message that client sent to message")
	fs.IntVar(&c.CycleTime, "CycleTime", 5, "The cycle time that client sent the message to master")
}
