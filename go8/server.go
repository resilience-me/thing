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

		// Only create a Conn if this is a client connection
		var sessionConn *Conn
		if datagram.Command&0x80 == 1 { // MSB is 1: Client connection
			sessionConn = &Conn{
				conn: conn,
				addr: remoteAddr,
			}
		}

		// Validate the datagram based on its type (client or server)
		if datagram.Command&0x80 == 0 { // Server session if MSB is 0
			if err := validateServerDatagram(buffer, datagram); err != nil {
				fmt.Printf("Error validating server datagram: %v\n", err)
				continue
			}
		} else { // Client session if MSB is 1
			errorMessage, err := validateClientDatagram(buffer, datagram)
			if err != nil {
				fmt.Printf("Error validating client datagram: %v\n", err)
				// Send an error response to the client
				sendErrorResponse(errorMessage, sessionConn)
				continue
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
