package option

import (
	"flag"
)

type ServerOption struct {
	ListenPort string
	HttpPort   string
	CycleTime  int
}

func NewServerOption() *ServerOption {
	return &ServerOption{}
}

func (c *ServerOption) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.ListenPort, "ListenPort", "7991", "The master url with ip and port")
	fs.StringVar(&c.HttpPort, "HttpPort", "6991", "The message that client sent to message")
	fs.IntVar(&c.CycleTime, "CycleTime", 5, "The cycle time that client sent the message to master")
}
