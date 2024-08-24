package main

import (
	"log"
	"fmt"
	"net"
	"ripple/pathfinding"
)

func main() {

	if err := config.InitConfig(); err != nil {
		log.Fatalf("Configuration failed: %v", err)
	}

	// Initialize the session manager
	sessionManager := NewSessionManager()

	// Initialize the path manager
	pathfinding.InitPathManager()

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

	fmt.Printf("Listening on port %d at server address %s.\n", config.Port, config.GetServerAddress())

	// Initialize the shutdown flag
	var shutdownFlag int32

	go shutdownHandler(conn, &shutdownFlag)

	// Start the server loop
	runServerLoop(conn, sessionManager, &shutdownFlag)

	sessionManager.wg.Wait()
	log.Println("All sessions and queues have been processed. Exiting.")
}
