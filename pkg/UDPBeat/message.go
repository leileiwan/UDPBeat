package UDPBeat

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/adler32"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// Message struct
type Message struct {
	ipSize   int32
	msgIP    string
	data     []byte
	checksum uint32
}

// NewMessage create a new message
func NewMessage(msgIP string, data string) *Message {
	ipSize := int32(len(msgIP))
	msg := &Message{
		ipSize: ipSize,
		msgIP:  msgIP,
		data:   []byte(data),
	}

	msg.checksum = msg.calcChecksum()
	return msg
}

// GetData get message data
func (msg *Message) GetData() []byte {
	return msg.data
}

// GetID get message ID
func (msg *Message) GetIP() string {
	return msg.msgIP
}

// Verify verify checksum
func (msg *Message) Verify() bool {
	return msg.checksum == msg.calcChecksum()
}

func (msg *Message) calcChecksum() uint32 {
	if msg == nil {
		return 0
	}

	data := new(bytes.Buffer)

	err := binary.Write(data, binary.LittleEndian, []byte(msg.msgIP))
	if err != nil {
		return 0
	}
	err = binary.Write(data, binary.LittleEndian, msg.data)
	if err != nil {
		return 0
	}

	checksum := adler32.Checksum(data.Bytes())
	return checksum
}

func (msg *Message) String() string {
	return fmt.Sprintf("IPSize=%d IP=%s DataLen=%d Checksum=%d", msg.ipSize, msg.GetIP(), len(msg.GetData()), msg.checksum)
}

// Encode from Message to []byte
func Encode(msg *Message) ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.LittleEndian, msg.ipSize)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buffer, binary.LittleEndian, []byte(msg.msgIP))
	if err != nil {
		return nil, err
	}
	err = binary.Write(buffer, binary.LittleEndian, msg.data)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buffer, binary.LittleEndian, msg.checksum)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(buffer.Bytes()))

	return buffer.Bytes(), nil
}

// Decode from []byte to Message
func Decode(data []byte) (*Message, error) {
	bufReader := bytes.NewReader(data)

	dataSize := len(data)
	// 读取ipsize
	var ipSize uint32
	err := binary.Read(bufReader, binary.LittleEndian, &ipSize)
	if err != nil {
		return nil, err
	}

	// 读取消息ID
	ipDataBuf := make([]byte, int(ipSize))
	err = binary.Read(bufReader, binary.LittleEndian, &ipDataBuf)
	if err != nil {
		return nil, err
	}
	msgIP := string(ipDataBuf)

	// 读取数据
	dataBufLength := dataSize - 4 - int(ipSize) - 4
	dataBuf := make([]byte, dataBufLength)
	err = binary.Read(bufReader, binary.LittleEndian, &dataBuf)
	if err != nil {
		return nil, err
	}

	// 检查checksum
	var checksum uint32
	err = binary.Read(bufReader, binary.LittleEndian, &checksum)
	if err != nil {
		return nil, err
	}

	message := &Message{}
	message.ipSize = int32(ipSize)
	message.msgIP = msgIP
	message.data = dataBuf
	message.checksum = checksum

	if message.Verify() {
		return message, nil
	}

	return nil, errors.New("checksum error")
}

func HostAddrCheck(addr string) bool {
	items := strings.Split(addr, ":")
	if items == nil || len(items) != 2 {
		return false
	}

	a := net.ParseIP(items[0])
	if a == nil {
		return false
	}

	match, err := regexp.MatchString("^[0-9]*$", items[1])
	if err != nil {
		return false
	}

	i, err := StringToInt64(items[1])
	if err != nil {
		return false
	}
	if i < 0 || i > 65535 {
		return false
	}

	if match == false {
		return false
	}

	return true
}

func StringToInt64(value string) (int64, error) {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}
