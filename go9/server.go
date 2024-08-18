package main

import (
	"fmt"
	"net"
	"ripple/config"
)

func runServerLoop(conn *net.UDPConn, sessionManager *SessionManager, ackManager *AckManager) {
	buffer := make([]byte, 393) // Combined buffer size (389 data + 4 ACK)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading from UDP connection: %v\n", err)
			continue
		}

		if n == 4 {
			// Handle client acknowledgment
			ackManager.mu.Lock()
			delete(ackManager.ackRegistry, string(buffer[:4]))
			ackManager.mu.Unlock()
			continue
		}

		if n != len(buffer) {
			fmt.Printf("Unexpected datagram size: received %d bytes, expected %d bytes\n", n, len(buffer))
			continue
		}

		fmt.Printf("Received %d bytes from %s\n", n, remoteAddr.String())

		// Extract the ACK part (first 4 bytes)
		ackBuffer := buffer[:4]
		
		// Extract the datagram part (remaining bytes)
		dataBuffer := buffer[4:]

		// Send an acknowledgment
		if err := udpr.SendAck(conn, remoteAddr, ackBuffer); err != nil {
			fmt.Printf("Failed to send ACK: %v\n", err)
			continue
		}

		// Parse the datagram
		datagram := parseDatagram(dataBuffer)

		// Validate the datagram
		if err := validateDatagram(dataBuffer, datagram); err != nil {
			fmt.Printf("Error validating  datagram: %v\n", err)
			continue
		}

		// Create a new session
		session := &Session{
			Datagram: datagram,
		}

		// If this is a client connection, associate the Conn with the session
		if datagram.Command&0x80 == 1 { // MSB is 1: Client connection
			session.Conn := &Client{
				UDPConn:     conn,
				Addr:        remoteAddr,
				AckManager:  ackManager,
			}
		}

		// Route the session through the SessionManager
		sessionManager.RouteSession(session)
	}
}

