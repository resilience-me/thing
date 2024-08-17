package main

import (
	"fmt"
	"net"
)

func runServerLoop(conn *net.UDPConn, transport *Transport, sessionManager *SessionManager) {
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

		// Parse the datagram
		datagram := parseDatagram(buffer)

		// Handle ACK datagrams separately
		if datagram.Command == 0x00 {
			ackKey := generateAckKey(datagram.Username, datagram.PeerUsername, datagram.PeerServerAddress, datagram.Counter)
			transport.RouteAck(ackKey)
			fmt.Println("ACK received and routed.")
			continue
		}

		var sessionConn *Conn

		// Determine if the datagram is from a server or client
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
			// Send ACK to server
			if err := SendServerAck(datagram); err != nil {
				fmt.Printf("Failed to send server ACK: %v\n", err)
			}
		} else { // Client session if MSB is 1
			errorMessage, err := validateClientDatagram(buffer, datagram)
			if err != nil {
				fmt.Printf("Error validating client datagram: %v\n", err)
				// Send an ACK with an error status and message to the client
				if err := SendClientAck(sessionConn, false, errorMessage); err != nil {
					fmt.Printf("Failed to send client error ACK: %v\n", err)
				}
				continue
			}
			// Send an ACK with a success status to the client
			if err := SendClientAck(sessionConn, true, ""); err != nil {
				fmt.Printf("Failed to send client success ACK: %v\n", err)
			}
		}

		// Create a new session with the appropriate Conn
		session := &Session{
			Datagram:  datagram,
			Conn:      sessionConn,
			Transport: transport, // Associate the Transport instance with the session
		}

		// Route the session through the SessionManager
		sessionManager.RouteSession(session)
	}
}
