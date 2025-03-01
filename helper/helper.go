package helper

import "os"

var (
	UserName     = ""
	Port         = 8080
	DebugMessage = make(chan string)
)

func GetOsHostName() string {
	name, err := os.Hostname()

	if err != nil {
		DebugMessage <- err.Error()
		panic(err)
	}

	return name
}
