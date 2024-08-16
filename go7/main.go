package main

import (
	"fmt"
	"net"
)

func main() {
	// Initialize the necessary components
	transport := NewTransport()
	sessionManager := NewSessionManager()

	// Set up the UDP server
	port := 2012
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Failed to listen on port %d: %v\n", port, err)
		return
	}
	defer conn.Close()

	// Start the server loop
	runServerLoop(conn, transport, sessionManager, port)
}
