package peer

import (
	"bufio"
	"fmt"
	"net"

	"github.com/sairash/p2p_chat_app/helper"
)

func Start(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		helper.MessageChan <- helper.DebugMessage(err.Error(), "Start")
		return
	}
	defer listener.Close()

	helper.MessageChan <- helper.DebugMessage(listener.Addr().String(), "Start")
	for {
		// Accept all the requests
		conn, err := listener.Accept()
		if err != nil {
			helper.MessageChan <- helper.DebugMessage(err.Error(), "Start")
			break
		}

		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	helper.MessageChan <- helper.DebugMessage(fmt.Sprintf("Connection opened by %s", conn.RemoteAddr()), "handleConnection")
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			helper.MessageChan <- helper.DebugMessage(fmt.Sprintf("Error reading from %s: %s \n", conn.RemoteAddr(), err.Error()), "handleConnection")
			return
		}

		helper.MessageChan <- helper.DebugMessage(string(message), "handleConnection")
		// fmt.Printf("[%s] Received: %s", conn.RemoteAddr(), message)
	}

}
