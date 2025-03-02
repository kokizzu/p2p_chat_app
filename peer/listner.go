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
		helper.DebugMessage <- err.Error()
		return
	}
	defer listener.Close()

	helper.DebugMessage <- fmt.Sprintf("%s \n", listener.Addr().String())
	for {
		// Accept all the requests
		conn, err := listener.Accept()
		if err != nil {
			helper.DebugMessage <- err.Error()
			break
		}

		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	helper.DebugMessage <- fmt.Sprintf("Connection opened by %s \n", conn.RemoteAddr())
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		text := scanner.Text()
		helper.DebugMessage <- fmt.Sprintf("%s \n", text)
		// fmt.Printf("Received from %s: %s\n", conn.RemoteAddr(), text)
	}

	if err := scanner.Err(); err != nil {
		helper.DebugMessage <- fmt.Sprintf("Error reading from %s: %v \n", conn.RemoteAddr(), err)
	}

	helper.DebugMessage <- fmt.Sprintf("Connection closed by %s \n", conn.RemoteAddr())
}
