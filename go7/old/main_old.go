// created with help of chat gpt, work in progress

package main

import (
	"fmt"
	"net"
)

type Datagram struct {
	Command           byte
	Username          string
	PeerUsername      string
	PeerServerAddress string
	Arguments         [256]byte
	Counter           uint32
	Signature         [32]byte
}

type Session struct {
	Datagram    *Datagram
	Conn        *Conn        // Pointer to Conn; can be nil
	AckRegistry *AckRegistry // Pointer to the AckRegistry
}

type Conn struct {
	conn *net.UDPConn  // Connection used to send/receive datagrams
	addr *net.UDPAddr  // Source address from where the datagram was received
}

// SessionManager manages sessions for different accounts
type SessionManager struct {
	sessionCh      chan *Session
	closedCh       chan string
	activeHandlers map[string]bool
	queues         map[string][]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessionCh:      make(chan *Session),
		closedCh:       make(chan string),
		activeHandlers: make(map[string]bool),
		queues:         make(map[string][]*Session),
	}
}

func (m *SessionManager) run() {
	for {
		select {
		case session := <-m.sessionCh:
			username := session.Datagram.Username
			if !m.activeHandlers[username] {
				m.activeHandlers[username] = true
				go m.handleSession(session)
			} else {
				m.queues[username] = append(m.queues[username], session)
			}

		case username := <-m.closedCh:
			if queue, exists := m.queues[username]; exists && len(queue) > 0 {
				nextSession := queue[0]
				m.queues[username] = queue[1:]
				go m.handleSession(nextSession)
			} else {
				delete(m.activeHandlers, username)
			}
		}
	}
}

func (m *SessionManager) handleSession(session *Session) {
	defer func() {
		m.closedCh <- session.Datagram.Username // Notify that the session is closed
	}()

	handleCommand(session)
}

func (m *SessionManager) addSession(session *Session) {
	m.sessionCh <- session
}

// CentralDispatcher routes incoming packets to the session manager
type CentralDispatcher struct {
	conn           *net.UDPConn
	transport      *Transport
	sessionManager *SessionManager
}

func NewCentralDispatcher(conn *net.UDPConn, sessionManager *SessionManager) *CentralDispatcher {
	transport := NewTransport() // No need to pass conn
	return &CentralDispatcher{
		conn:           conn,
		transport:      transport,
		sessionManager: sessionManager,
	}
}

func (cd *CentralDispatcher) ListenAndServe() {
	buffer := make([]byte, 389)
	for {
		n, addr, err := cd.conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error receiving UDP packet:", err)
			continue
		}

		packet := buffer[:n]
		messageType := packet[0] // Assuming the first byte indicates the message type

		switch messageType {
		case 0x00: // Specific case for 0x00 (ACK)
			ack := deserializeAck(packet)
			cd.transport.RouteAck(ack) // Use the embedded AckRegistry methods
		default: // All other cases (Datagrams)
			var conn *Conn = nil
			if messageType&0x80 == 0 {
				// MSB is 0: Use the existing conn and addr
				conn = &Conn{conn: cd.conn, addr: addr}
			}

			datagram := deserializeDatagram(packet)
			session := &Session{
				Datagram:    datagram,
				Conn:        conn,
				AckRegistry: cd.transport.AckRegistry, // Set the AckRegistry from Transport
			}
			cd.sessionManager.addSession(session)
		}
	}
}

func handleCommand(session *Session) {
	fmt.Printf("Handling command %d for %s -> %s\n", session.Datagram.Command, session.Datagram.Username, session.Datagram.PeerUsername)

	// Example of sending a response back to the source using transport logic
	if session.Conn != nil {
		responseData := []byte("Response Data")
		err := session.Conn.conn.WriteToUDP(responseData, session.Conn.addr)
		if err != nil {
			fmt.Printf("Failed to send response: %v\n", err)
		}
	}
}

// Utility functions to serialize and deserialize Datagram and Ack structures
func deserializeAck(data []byte) *Ack {
	// Simplified deserialization logic
	return &Ack{}
}

func deserializeDatagram(data []byte) *Datagram {
	// Simplified deserialization logic
	return &Datagram{}
}

func serializeDatagram(dg *Datagram) []byte {
	// Simplified serialization logic
	return []byte{}
}

func main() {
	// Listen on all interfaces, port 2012
	addr := net.UDPAddr{Port: 2012, IP: net.ParseIP("0.0.0.0")}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Failed to bind to UDP port: %v\n", err)
		return
	}
	defer conn.Close()

	// Initialize the session manager and central dispatcher
	sessionManager := NewSessionManager()
	dispatcher := NewCentralDispatcher(conn, sessionManager)

	// Run the session manager and central dispatcher
	go sessionManager.run()
	go dispatcher.ListenAndServe()

	// Block the main goroutine
	select {}
}
