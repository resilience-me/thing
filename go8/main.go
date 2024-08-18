package main

import (
	"fmt"
	"net"
)

func main() {
	// Initialize the session manager and ackManager
	sessionManager := NewSessionManager()
	ackManager := NewAckManager()

	// Set up the UDP server
	addr := net.UDPAddr{
		Port: Port,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Failed to listen on port %d: %v\n", Port, err)
		return
	}
	defer conn.Close()

	// Start the server loop
	runServerLoop(conn, sessionManager, ackManager)
}
