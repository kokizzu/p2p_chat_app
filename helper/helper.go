package helper

import (
	"fmt"
	"net"
	"os"
)

type HInfo struct {
	Ip   string
	Port string
	Name string
	Conn net.Conn
}

var (
	IP          = ""
	IPPORT      = ""
	UserName    = ""
	Local       = false
	Debug       = false
	Port        = 8080
	MessageChan = make(chan DisplayMessage)

	ConnectedHosts = map[string]HInfo{}
)

const (
	DebugType = iota
	Self
	Peer
)

type Message struct {
	Text string
	Name string
	IP   string
}
type DisplayMessage struct {
	TypeOfMessage uint
	Message
}

func GetNameFromIP(addr string) string {
	if in_val, in_has := ConnectedHosts[addr]; in_has {
		return in_val.Name
	}
	return "Peer"
}

func GetOsHostName() string {
	name, err := os.Hostname()

	if err != nil {
		MessageChan <- DebugMessage(err.Error(), "GetHostName")
		panic(err)
	}

	return name
}

func DebugMessage(message string, from string) DisplayMessage {
	return DisplayMessage{
		TypeOfMessage: DebugType,
		Message: Message{
			Text: fmt.Sprintf("%s \n", message),
			Name: from,
		},
	}

}
