// written together with chat gpt, work in progress

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

	buffer := make([]byte, 1024) // Buffer to hold incoming datagrams

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading from UDP connection: %v\n", err)
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
		if buffer[0]&0x80 == 0 { // MSB is 0: Server connection
			sessionConn = &Conn{
				conn: conn,
				addr: remoteAddr,
			}
		} else { // MSB is 1: Client connection
			sessionConn = nil
		}

		// Create a new session with the appropriate Conn
		session := &Session{
			Datagram:  datagram,
			Conn:      sessionConn,
			Transport: transport,  // Associate the Transport instance with the session
		}

		// Route the session through the SessionManager
		sessionManager.RouteSession(session)
	}
}
