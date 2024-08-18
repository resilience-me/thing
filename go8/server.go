package main

import (
	"fmt"
	"net"
)

func runServerLoop(conn *net.UDPConn, sessionManager *SessionManager) {
	buffer := make([]byte, 393) // Combined buffer size (389 data + 4 ACK)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading from UDP connection: %v\n", err)
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
		dataBuffer := buffer[4:n]

		// Create a Conn object for the acknowledgment and potential session
		remoteConn := &Conn{
			UDPConn: conn,
			addr:    remoteAddr,
		}

		// Send an acknowledgment back to the client as soon as possible
		if err := Ack(remoteConn, ackBuffer); err != nil {
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
			session.Conn = remoteConn
		}

		// Route the session through the SessionManager
		sessionManager.RouteSession(session)
	}
}

