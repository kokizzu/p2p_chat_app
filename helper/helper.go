package helper

import "os"

var (
	UserName     = ""
	Debug        = false
	Local        = false
	Port         = 8080
	DebugMessage = make(chan string)
)

type Message struct {
	Text string
	Name string
}
type DisplayMessage struct {
	Self bool
	Message
}

func GetOsHostName() string {
	name, err := os.Hostname()

	if err != nil {
		DebugMessage <- err.Error()
		panic(err)
	}

	return name
}
