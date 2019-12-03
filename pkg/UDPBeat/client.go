package UDPBeat

import (
	"fmt"
	"net"
	"time"
)

// SocketClinet struct
type SocketClient struct {
	serverAddr string
	stopCh     chan error
	msg        *Message
	cycleTime  time.Duration
}

func (sc *SocketClient) Serv() {
	sc.sentHandler()
	for {
		select {
		case <-sc.stopCh:
			fmt.Println("The client end...")
			return
		}
	}
}

func NewSockerClient(serverAddr, data string, cycleTime time.Duration) (*SocketClient, error) {
	localIP, err := GetInternal()
	if err != nil {
		return nil, err
	}
	msg := NewMessage(localIP, data)
	return &SocketClient{
		serverAddr: serverAddr,
		stopCh:     make(chan error),
		msg:        msg,
		cycleTime:  cycleTime,
	}, nil
}

// tell the master I'am alive
func (sc *SocketClient) Beat() error {
	conn, err := net.DialTimeout("udp", sc.serverAddr, time.Second*5)
	if err != nil {
		return err
	}
	defer conn.Close()
	msgBytes, err := Encode(sc.msg)
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte(msgBytes))
	return err
}

func (sc *SocketClient) sentHandler() {
	go func() {
		for {
			err := sc.Beat()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Send the message %v\n", sc.msg)
			time.Sleep(sc.cycleTime)
		}
	}()

}

func GetInternal() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	var ip string
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.To4().String()
			}
		}
	}
	return ip, err
}

//add for test
func (sc *SocketClient) SetRecycleTime(time time.Duration) {
	sc.cycleTime = time
}
