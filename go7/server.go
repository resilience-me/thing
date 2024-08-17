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

		var sessionConn *Conn = nil

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
	        } else { // Client session if MSB is 1
	            errorMessage, err := validateClientDatagram(buffer, datagram)
	            if err != nil {
	                fmt.Printf("Error validating client datagram: %v\n", err)
	                // Optionally, send an error response to the client
	                sendErrorResponse(errorMessage, sessionConn)
	                continue
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
