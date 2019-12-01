package UDPBeat

import (
	"fmt"
	"net"
)

// SocketService struct
type SocketService struct {
	onConnect func(chan Message)
	laddr     string
	listener  net.PacketConn
	stopCh    chan error
}

// NewSocketService create a new socket service
func NewSocketService(laddr string) (*SocketService, error) {

	// l, err := net.Listen("udp", laddr)
	l, err := net.ListenPacket("udp", laddr)

	if err != nil {
		return nil, err
	}
	s := &SocketService{
		stopCh:   make(chan error),
		listener: l,
	}

	return s, nil
}

// RegConnectHandler register connect handler
func (s *SocketService) RegConnectHandler(handler func(chan Message)) {
	s.onConnect = handler
}

// Serv Start socket service
func (s *SocketService) Serv() {
	defer func() {
		s.listener.Close()
	}()

	ch, err := s.acceptHandler()
	if err != nil {
		fmt.Println(err)
	}

	s.onConnect(ch)

	for {
		select {

		case <-s.stopCh:
			fmt.Println("The server end...")
			return
		}
	}
}

func (s *SocketService) acceptHandler() (chan Message, error) {
	ch := make(chan Message)
	buf := make([]byte, 1000)
	go func() {
		for {
			//recive bytes
			n, addr, err := s.listener.ReadFrom(buf)
			if err != nil {
				fmt.Println(err)
				continue
			}
			//decode bytes to Message
			msg, err := Decode(buf[0:n])
			fmt.Println("Receive from", addr, msg)
			if err != nil {
				fmt.Println(err)
				continue
			}

			ch <- *msg

		}
	}()

	return ch, nil
}
