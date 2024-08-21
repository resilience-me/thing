package main

import (
	"log"
	"fmt"
	"net"
)

func main() {

	if err := config.InitConfig(); err != nil {
		log.Fatalf("Configuration failed: %v", err)
	}

	// Initialize the session manager
	sessionManager := NewSessionManager()

	// Set up the UDP server
	addr := net.UDPAddr{
		Port: config.Port,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Failed to listen on port %d: %v\n", Port, err)
		return
	}
	defer conn.Close()

	fmt.Printf("Server is running at address: %s\n", config.GetServerAddress())
	fmt.Println("Listening on port 2012...")

	// Start the server loop
	runServerLoop(conn, sessionManager)
}
