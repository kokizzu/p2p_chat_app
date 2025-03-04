package peer

import (
	"fmt"
	"net"
	"strings"
	"time"

	reuseport "github.com/kavu/go_reuseport"
	"github.com/sairash/p2p_chat_app/helper"
)

const (
	portBroadcast       = 7715
	bufferSize          = 1024
	messageEncoderSplit = "||"
)

func StartDiscovery() {
	go func() {
		// Broadcast presence periodically
		for {
			broadcastToDiscover()
			time.Sleep(time.Second * 5)
		}
	}()

	// Create persistent listener once
	pc, err := reuseport.ListenPacket("udp4", fmt.Sprintf("0.0.0.0:%d", portBroadcast))
	if err != nil {
		helper.MessageChan <- helper.DebugMessage(fmt.Sprintf("Failed to create listener: %v", err), "StartDiscovery")
		return
	}
	defer pc.Close()

	helper.MessageChan <- helper.DebugMessage("Listening to Broadcast", "StartDiscovery")

	// Continuous read loop
	for {
		listenForBroadcast(pc)
	}
}

func listenForBroadcast(pc net.PacketConn) {
	buffer := make([]byte, bufferSize)

	// Set read deadline to prevent permanent blocking
	pc.SetReadDeadline(time.Now().Add(2 * time.Millisecond))

	bytesRead, addr, err := pc.ReadFrom(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// Timeout is expected, just return and try again
			return
		}
		helper.MessageChan <- helper.DebugMessage(fmt.Sprintf("Error reading from UDP: %v", err), "listenForBroadcast")
		return
	}

	remoteAddr, ok := addr.(*net.UDPAddr)
	if !ok {
		helper.MessageChan <- helper.DebugMessage("Invalid remote address type", "listenForBroadcast")
		return
	}

	messageAttributes := strings.Split(string(buffer[:bytesRead]), messageEncoderSplit)
	if len(messageAttributes) != 2 {
		helper.MessageChan <- helper.DebugMessage("Invalid message format", "listenForBroadcast")
		return
	}

	ip := remoteAddr.IP.String()
	peerAddress := ip + ":" + messageAttributes[0]

	// Add peer to connected hosts
	if peerAddress == helper.IPPORT {
		return
	}
	if _, has := helper.ConnectedHosts[peerAddress]; !has {
		conn, err := net.Dial("tcp", peerAddress)
		if err != nil {
			helper.MessageChan <- helper.DebugMessage(fmt.Sprintf("Couldn't connect to the address %s because of %s", peerAddress, err.Error()), "ConnnectHost")
			return
		}
		helper.ConnectedHosts[peerAddress] = helper.HInfo{
			Ip:   ip,
			Port: messageAttributes[0],
			Name: messageAttributes[1],
			Conn: conn,
		}
		helper.MessageChan <- helper.DebugMessage(fmt.Sprintf("Discovered peer: %s", peerAddress), "listenForBroadcast")
	}
}

func broadcastToDiscover() {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", net.IPv4bcast.String(), portBroadcast))
	if err != nil {
		helper.MessageChan <- helper.DebugMessage(fmt.Sprintf("Error resolving broadcast address: %v", err), "broadcastToDiscover")
		return
	}

	udpConnection, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		helper.MessageChan <- helper.DebugMessage(fmt.Sprintf("Error creating UDP connection: %v", err), "broadcastToDiscover")
		return
	}
	defer udpConnection.Close()

	buffer := []byte(fmt.Sprintf("%d%s%s", helper.Port, messageEncoderSplit, helper.UserName))
	_, err = udpConnection.Write(buffer)

	if err != nil {
		helper.MessageChan <- helper.DebugMessage(fmt.Sprintf("Error while broadcasting: %v", err), "broadcastToDiscover")
	}
}

// Getting local address with dns seems much easier
func GetHostIPAddressV4() string {
	conn, err := net.Dial("udp", "1.1.1.1:53")

	if err != nil {
		helper.MessageChan <- helper.DebugMessage(err.Error(), "GetHostIP")
		panic(err)
	}

	helper.MessageChan <- helper.DebugMessage(strings.Split(conn.LocalAddr().String(), ":")[0], "GetHostIP")
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}
