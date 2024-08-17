package main

import (
	"fmt"
	"net"
)

func runServerLoop(conn *net.UDPConn, transport *Transport, sessionManager *SessionManager, port int) {
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
		if err := ValidateDatagram(buffer[:n], datagram); err != nil {
			fmt.Printf("Datagram validation failed: %v\n", err)
			continue
		}

		var sessionConn *Conn = nil

		// Determine if the datagram is from a server or client
		if datagram.Command&0x80 == 1 { // MSB is 1: Client connection
			sessionConn = &Conn{
				conn: conn,
				addr: remoteAddr,
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

		// Determine the ACK address
		var ackAddr *net.UDPAddr
		if datagram.Command&0x80 == 0 { // Server connection
			ackAddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", datagram.PeerServerAddress, port))
			if err != nil {
				fmt.Printf("Failed to resolve server address: %v\n", err)
				continue
			}
		} else { // Client connection
			ackAddr = remoteAddr
		}

		// Generate, sign, and serialize the ACK datagram
		ackData := generateAndSignAckDatagram(datagram)

		// Send the ACK back to the determined address
		err = SendAck(ackData, ackAddr)
		if err != nil {
			fmt.Printf("Failed to send ACK to %s: %v\n", ackAddr.String(), err)
		} else {
			fmt.Printf("Sent ACK to %s\n", ackAddr.String())
		}
	}
}