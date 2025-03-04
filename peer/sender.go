package peer

import (
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
	for k, v := range helper.ConnectedHosts {
		_, err := v.Conn.Write([]byte(message))

		if err != nil {
			helper.MessageChan <- helper.DebugMessage(fmt.Sprintf("Couldn't send mesasge to %s because of %s", k, err.Error()), "Send")
		}
	}
}
