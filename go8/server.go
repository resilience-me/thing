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

		// Validate the datagram
		if err := ValidateDatagram(dg); err != nil {
			fmt.Printf("Error validating datagram: %v\n", err)
			continue
		}

		// Validate the datagram
		alreadyInQueue, err := ValidateAndIncrementCounter(dg)
		if err != nil {
			fmt.Printf("Error validating counter: %v\n", err)
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

		// Send ack for datagram based on its type (client or server)
		if datagram.Command&0x80 == 0 { // Server session if MSB is 0
			// Send ACK to server
			if err := SendServerAck(datagram); err != nil {
				fmt.Printf("Failed to send server ACK: %v\n", err)
			}
		} else { // Client session if MSB is 1
			// Send an ACK with a success status to the client
			if err := SendClientAck(sessionConn); err != nil {
				fmt.Printf("Failed to send client ACK: %v\n", err)
			}
		}

		if alreadyInQueue {
			continue
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
