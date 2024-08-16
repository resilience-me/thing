package main

import (
	"fmt"
	"net"
)

func main() {
	// Create the Transport instance
	transport := NewTransport()

	// Create the SessionManager
	sessionManager := NewSessionManager()

	// Set up the UDP server
	port := 2012
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Failed to listen on port %d: %v\n", port, err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 389) // Buffer sized according to datagram size

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading from UDP connection: %v\n", err)
			continue
		}

		// Ensure that exactly 389 bytes were received
		if n != 389 {
			fmt.Printf("Unexpected datagram size: received %d bytes, expected 389 bytes\n", n)
			continue
		}

		// Process the received datagram
		fmt.Printf("Received %d bytes from %s\n", n, remoteAddr.String())

		// Parse the datagram
		datagram := parseDatagram(buffer[:n])

		// Determine if the MSB of the first byte (Command) is 1 or 0
		// If MSB is 1, it's a client connection, so set Conn to nil
		// If MSB is 0, it's a server connection, so include the Conn with the address
		var sessionConn *Conn
		var ackAddr string

		if buffer[0]&0x80 == 0 { // MSB is 0: Server connection
			sessionConn = &Conn{
				conn: conn,
				addr: remoteAddr,
			}
			// Send ACK to the server address specified in the datagram
			ackAddr = fmt.Sprintf("%s:%d", datagram.PeerServerAddress, port)
		} else { // MSB is 1: Client connection
			sessionConn = nil
			// Send ACK back to the client address from which the datagram was received
			ackAddr = remoteAddr.String()
		}

		// Create a new session with the appropriate Conn
		session := &Session{
			Datagram:  datagram,
			Conn:      sessionConn,
			Transport: transport,  // Associate the Transport instance with the session
		}

		// Route the session through the SessionManager
		sessionManager.RouteSession(session)

		// Determine if the datagram is an ACK (e.g., Command 0x00)
		if datagram.Command == 0x00 {
			fmt.Println("Received an ACK, no response necessary.")
			continue
		}

		// Otherwise, send an ACK back to the determined address
		ack := NewAck(datagram)
		ackData := serializeDatagram(ack)
		err = SendAck(ackData, ackAddr)
		if err != nil {
			fmt.Printf("Failed to send ACK to %s: %v\n", ackAddr, err)
		} else {
			fmt.Printf("Sent ACK to %s\n", ackAddr)
		}
	}
}
