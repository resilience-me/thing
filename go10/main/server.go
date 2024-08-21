package main

import (
	"fmt"
	"net"
	"ripple/auth"
	"ripple/comm"
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
		dataBuffer := buffer[4:]

		// Send an acknowledgment
		if err := comm.SendAck(conn, remoteAddr, ackBuffer); err != nil {
			fmt.Printf("Failed to send ACK: %v\n", err)
			continue
		}

		// Parse the datagram
		datagram := DeserializeDatagram(dataBuffer)

		// Validate the datagram
		if err := auth.ValidateDatagram(dataBuffer, datagram); err != nil {
			fmt.Printf("Error validating  datagram: %v\n", err)
			continue
		}

		// Create a new session
		session := &Session{
			Datagram: datagram,
		}

		// If this is a client connection, associate the Conn with the session
		if datagram.Command&0x80 != 0 { // MSB is 1: Client connection
			// Create a Conn object for the session
			session.Conn = &Conn{
				UDPConn: conn,
				addr:    remoteAddr,
			}
		}

		// Route the session through the SessionManager
		sessionManager.RouteSession(session)
	}
}
