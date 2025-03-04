package peer

import (
	"encoding/json"
	"fmt"

	"github.com/sairash/p2p_chat_app/helper"
)

// func ConnnectHost(address string) {
// 	_, has := helper.ConnectedHosts[address]

// 	if has {
// 		helper.MessageChan <- helper.DebugMessage(fmt.Sprintf("%s address already connected.", address), "ConnnectHost")
// 		return
// 	}

// 	helper.MessageChan <- helper.DebugMessage(fmt.Sprintf("Connecting to address %s", address), "ConnnectHost")

// 	helper.ConnectedHosts[address] = conn
// }

func Send(message string) {
	message_struct, err := json.Marshal(helper.Message{Text: message, Name: "", IP: helper.IPPORT})
	if err != nil {
		helper.MessageChan <- helper.DebugMessage("Couldn't send message due to json.", "Send")
		return
	}
	message_struct = append(message_struct, '\n')
	for k, v := range helper.ConnectedHosts {
		_, err := v.Conn.Write(message_struct)

		if err != nil {
			helper.MessageChan <- helper.DebugMessage(fmt.Sprintf("Couldn't send mesasge to %s because of %s", k, err.Error()), "Send")
		}
	}
}
