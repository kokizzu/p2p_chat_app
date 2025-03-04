package peer

import (
	"bufio"
	"encoding/json"
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

		message_struct := helper.Message{}

		err = json.Unmarshal(message, &message_struct)

		if err != nil {
			helper.MessageChan <- helper.DebugMessage("Couldn't Structure Message.", "handleConnection")
		}

		helper.MessageChan <- helper.DisplayMessage{
			Message: helper.Message{
				Text: message_struct.Text,
				Name: helper.GetNameFromIP(message_struct.IP),
			},
			TypeOfMessage: helper.Peer,
		}
		// fmt.Printf("[%s] Received: %s", conn.RemoteAddr(), message)
	}

}
