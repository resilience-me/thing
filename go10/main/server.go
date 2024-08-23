package main

import (
	"log"
	"net"
	"sync/atomic"
	"ripple/auth"
	"ripple/comm"
	"ripple/types"
	"ripple/pathfinding"
)

// runServerLoop runs the main server loop, processing incoming datagrams
func runServerLoop(conn *net.UDPConn, sessionManager *SessionManager, pathManager *pathfinding.PathManager, shutdownFlag *int32) {
	buffer := make([]byte, 393) // Combined buffer size (389 data + 4 ACK)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			// Check if the error is because of a shutdown (e.g., conn.Close() was called)
			if atomic.LoadInt32(shutdownFlag) != 0 {
				log.Println("Server is shutting down...")
				return
			}
			// Handle other errors (unexpected issues)
			log.Printf("Error reading from UDP connection: %v", err)
			continue
		}

		if n != len(buffer) {
			log.Printf("Unexpected datagram size: received %d bytes, expected %d bytes", n, len(buffer))
			continue
		}

		log.Printf("Received %d bytes from %s", n, remoteAddr.String())

		// Extract the ACK part (first 4 bytes)
		ackBuffer := buffer[:4]
		
		// Extract the datagram part (remaining bytes)
		dataBuffer := buffer[4:]

		// Send an acknowledgment
		if err := comm.SendAck(conn, remoteAddr, ackBuffer); err != nil {
			log.Printf("Failed to send ACK: %v", err)
			continue
		}

		// Parse the datagram
		datagram := types.DeserializeDatagram(dataBuffer)

		// Validate the datagram
		if err := auth.ValidateDatagram(dataBuffer, datagram); err != nil {
			log.Printf("Error validating datagram: %v", err)
			continue
		}

		// Create a new session
		session := &Session{
			Datagram: datagram,
			PathManager: pathManager,
		}

		// If this is a client connection, associate the address with the session
		if datagram.Command&0x80 == 0 { // MSB is 0: Client connection
			// Add the client address to the session object
			session.Addr = remoteAddr,
		}

		// Route the session through the SessionManager
		sessionManager.RouteSession(session)
	}
}
