package option

import (
	"flag"
	"time"
)

type ServerOption struct {
	ListenPort   string
	HttpPort     string
	CycleTime    time.Duration
	LevelInitHP  int
	LevelAliveHP int
	LevelFullHP  int
}

func NewServerOption() *ServerOption {
	return &ServerOption{}
}

func (c *ServerOption) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.ListenPort, "ListenPort", "7991", "The master url with ip and port")
	fs.StringVar(&c.HttpPort, "HttpPort", "6991", "The message that client sent to message")
	fs.DurationVar(&c.CycleTime, "CycleTime", 5*time.Second, "The cycle time that client sent the message to master")
	fs.IntVar(&c.LevelInitHP, "LevelInitHP", 0, "The cycle time that client sent the message to master")
	fs.IntVar(&c.LevelAliveHP, "LevelAliveHP", 1, "The cycle time that client sent the message to master")
	fs.IntVar(&c.LevelFullHP, "LevelFullHP", 5, "The cycle time that client sent the message to master")

}
