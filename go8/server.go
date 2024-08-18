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

		// Create a Conn object for the acknowledgment and potential session
		remoteConn := &Conn{
			UDPConn: conn,
			addr:    remoteAddr,
		}

		// Send an acknowledgment back to the client as soon as possible
		if err := SendAck(remoteConn); err != nil {
			fmt.Printf("Failed to send ACK: %v\n", err)
			continue
		}

		// Parse the datagram
		datagram := parseDatagram(buffer)

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
				// Send an error response to the client in a separate goroutine
				go SendErrorResponse(errorMessage, remoteConn)
				continue
			}
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

