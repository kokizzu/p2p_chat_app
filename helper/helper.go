package helper

import (
	"fmt"
	"os"
)

type HInfo struct {
	Ip   string
	Port string
	Name string
}

var (
	IP          = ""
	IPPORT      = ""
	UserName    = ""
	Local       = false
	Port        = 8080
	MessageChan = make(chan DisplayMessage)

	ConnectedHosts = map[string]HInfo{}
)

const (
	Debug = iota
	Self
	Peer
)

type Message struct {
	Text string
	Name string
}
type DisplayMessage struct {
	TypeOfMessage uint
	Message
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
		TypeOfMessage: Debug,
		Message: Message{
			Text: fmt.Sprintf("%s \n", message),
			Name: from,
		},
	}

}
