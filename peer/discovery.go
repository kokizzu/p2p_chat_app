package peer

import (
	"net"
	"strings"

	"github.com/sairash/p2p_chat_app/helper"
)

// Getting local address with dns seems much easier
func GetHostIPAddress() string {
	conn, err := net.Dial("udp", "1.1.1.1:53")

	if err != nil {
		helper.DebugMessage <- err.Error()
		panic(err)
	}

	helper.DebugMessage <- strings.Split(conn.LocalAddr().String(), ":")[0] + " \n"
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}
