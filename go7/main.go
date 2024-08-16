package main

import (
	"fmt"
	"net"
)

func main() {
	// Initialize the necessary components
	transport := NewTransport()
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

	// Start the server loop with everything inlined
	runServerLoop(conn, transport, sessionManager, port)
}

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

		datagram := parseDatagram(buffer[:n])

		if datagram.Command == 0x00 {
			ackKey := generateAckKey(datagram.Username, datagram.PeerUsername, datagram.PeerServerAddress, datagram.Counter)
			transport.RouteAck(ackKey)
			fmt.Println("ACK received and routed.")
			continue
		}

		var sessionConn *Conn
		var ackAddr *net.UDPAddr

		if datagram.Command&0x80 == 0 { // MSB is 0: Server connection
			sessionConn = &Conn{
				conn: conn,
				addr: remoteAddr,
			}
			ackAddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", datagram.PeerServerAddress, port))
			if err != nil {
				fmt.Printf("Failed to resolve server address: %v\n", err)
				continue
			}
		} else { // MSB is 1: Client connection
			sessionConn = nil
			ackAddr = remoteAddr
		}

		session := &Session{
			Datagram:  datagram,
			Conn:      sessionConn,
			Transport: transport, // Associate the Transport instance with the session
		}

		sessionManager.RouteSession(session)

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
