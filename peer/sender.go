package peer

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/sairash/p2p_chat_app/helper"
)

func Send(address string, message string, name string) {
	message_struct, err := json.Marshal(helper.Message{Text: message, Name: name})
	if err != nil {
		helper.DebugMessage <- "Couldn't send message due to json. \n"
		return
	}
	conn, err := net.Dial("tcp", address)
	if err != nil {

		helper.DebugMessage <- fmt.Sprintf("Couldn't send mesasge to %s because of %s \n", address, err.Error())
		return
	}

	_, err = conn.Write(message_struct)

	if err != nil {
		helper.DebugMessage <- fmt.Sprintf("Couldn't send mesasge to %s because of %s \n", address, err.Error())
	}
}
