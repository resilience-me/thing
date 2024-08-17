package main

import (
	"fmt"
	"net"
)

func runServerLoop(conn *net.UDPConn, sessionManager *SessionManager) {
	buffer := make([]byte, 389) // Buffer sized according to datagram size

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

		// Send an acknowledgment back to the client as soon as possible
		if err := SendAck(conn, remoteAddr); err != nil {
			fmt.Printf("Failed to send ACK: %v\n", err)
			continue
		}

		// Parse the datagram
		datagram := parseDatagram(buffer)

		// Application layer validation
		if err := ValidateDatagram(buffer, datagram); err != nil {
			fmt.Printf("Error validating datagram: %v\n", err)
			continue
		}

		// Validate and increment counter (for business logic, not transport)
		if _, err := ValidateAndIncrementCounter(datagram); err != nil {
			fmt.Printf("Error validating counter: %v\n", err)
			continue
		}

		// Only create a Conn if this is a client connection
		var sessionConn *Conn
		if datagram.Command&0x80 == 1 { // MSB is 1: Client connection
			sessionConn = &Conn{
				conn: conn,
				addr: remoteAddr,
			}
		}

		// Create a new session with the appropriate Conn (which could be nil)
		session := &Session{
			Datagram: datagram,
			Conn:     sessionConn,
		}

		// Route the session through the SessionManager
		sessionManager.RouteSession(session)
	}
}
