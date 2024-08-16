// work in progress, chat gpt helping out

package main

import (
	"fmt"
	"net"
	"sync"
	"time"
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

type Conn struct {
	conn *net.UDPConn  // Connection used to send/receive datagrams
	addr *net.UDPAddr  // Source address from where the datagram was received
}

type Session struct {
	Datagram     Datagram
	conn         Conn          // Contains both UDPConn and UDPAddr
	ackRegistry  *AckRegistry  // Pointer to the AckRegistry
}

type Ack struct {
	Username          string
	PeerUsername      string
	PeerServerAddress string
	Counter           uint32
}

// SyncManager manages synchronization for different accounts
type SyncManager struct {
	mu       sync.Mutex
	syncMap  map[string]*sync.Mutex
}

func NewSyncManager() *SyncManager {
	return &SyncManager{
		syncMap: make(map[string]*sync.Mutex),
	}
}

func (sm *SyncManager) getMutex(key string) *sync.Mutex {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.syncMap[key]; !exists {
		sm.syncMap[key] = &sync.Mutex{}
	}
	return sm.syncMap[key]
}

// AckRegistry manages ACKs for different accounts
type AckRegistry struct {
	mu          sync.Mutex
	waitingAcks map[string]chan Ack
}

func NewAckRegistry() *AckRegistry {
	return &AckRegistry{
		waitingAcks: make(map[string]chan Ack),
	}
}

// generateBaseKey creates a common key based on username, peer username, and peer server address
func generateBaseKey(username, peerUsername, peerServerAddress string) string {
	return fmt.Sprintf("%s-%s-%s", username, peerUsername, peerServerAddress)
}

// generateAckKey creates a unique key for ACKs based on the base key and counter
func generateAckKey(username, peerUsername, peerServerAddress string, counter uint32) string {
	return fmt.Sprintf("%s-%d", generateBaseKey(username, peerUsername, peerServerAddress), counter)
}

// generateCommandKey creates a unique key based on username for command handling
func generateCommandKey(username string) string {
	return username
}

func (ar *AckRegistry) registerAck(ack Ack) chan Ack {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	key := generateAckKey(ack.Username, ack.PeerUsername, ack.PeerServerAddress, ack.Counter)
	ch := make(chan Ack)
	ar.waitingAcks[key] = ch
	return ch
}

func (ar *AckRegistry) routeAck(ack Ack) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	key := generateAckKey(ack.Username, ack.PeerUsername, ack.PeerServerAddress, ack.Counter)
	if ch, exists := ar.waitingAcks[key]; exists {
		ch <- ack
		close(ch) // Close the channel after sending the ACK to avoid leaks
		delete(ar.waitingAcks, key)
	}
}

func (ar *AckRegistry) cleanupAck(ack Ack) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	key := generateAckKey(ack.Username, ack.PeerUsername, ack.PeerServerAddress, ack.Counter)
	delete(ar.waitingAcks, key)
}

type CentralDispatcher struct {
	conn         *net.UDPConn
	ackRegistry  *AckRegistry
	syncManager  *SyncManager
}

func NewCentralDispatcher(conn *net.UDPConn, ackRegistry *AckRegistry, syncManager *SyncManager) *CentralDispatcher {
	return &CentralDispatcher{
		conn:        conn,
		ackRegistry: ackRegistry,
		syncManager: syncManager,
	}
}

func (cd *CentralDispatcher) ListenAndServe() {
	buffer := make([]byte, 1024)
	for {
		n, addr, err := cd.conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error receiving UDP packet:", err)
			continue
		}

		packet := buffer[:n]
		messageType := packet[0] // Assuming the first byte indicates the message type

		switch messageType & 0x80 { // Check the MSB
		case 0x00: // Server connection (MSB is 0)
			ack := deserializeAck(packet)
			cd.ackRegistry.routeAck(ack)
		case 0x80: // Client connection (MSB is 1)
			datagram := deserializeDatagram(packet)
			session := createSession(datagram, cd.conn, addr, cd.ackRegistry) // Include sourceAddr
			cd.routeToCommandHandler(session, datagram, addr)
		}
	}
}

func (cd *CentralDispatcher) routeToCommandHandler(session *Session, datagram Datagram, addr *net.UDPAddr) {
	// Create a key for the mutex based on username for command handling
	mutex := cd.syncManager.getMutex(generateCommandKey(datagram.Username))
	mutex.Lock()
	defer mutex.Unlock()

	// Process the command synchronously
	handleCommand(session, datagram, addr)
}

func handleCommand(session *Session, datagram Datagram, addr *net.UDPAddr) {
	fmt.Printf("Handling command %d for %s -> %s\n", datagram.Command, datagram.Username, datagram.PeerUsername)

	// Example of sending a response back to the source
	responseData := []byte("Response Data")
	err := session.conn.WriteToUDP(responseData, addr)
	if err != nil {
		fmt.Printf("Failed to send response: %v\n", err)
	}
}

func sendPacketWithRetry(session *Session, packet []byte, maxRetries int) error {
	retries := 0
	delay := 1 * time.Second

	serverAddress := session.Datagram.PeerServerAddress
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", serverAddress, 2012))
	if err != nil {
		return fmt.Errorf("failed to resolve server address '%s': %w", serverAddress, err)
	}

	ack := Ack{
		Username:          session.Datagram.Username,
		PeerUsername:      session.Datagram.PeerUsername,
		PeerServerAddress: session.Datagram.PeerServerAddress,
		Counter:           session.Datagram.Counter,
	}

	ackChan := session.ackRegistry.registerAck(ack)

	for retries < maxRetries {
		if _, err := session.conn.WriteToUDP(packet, addr); err != nil {
			return fmt.Errorf("failed to send data to server '%s': %w", serverAddress, err)
		}

		select {
		case receivedAck := <-ackChan:
			if receivedAck.Counter == session.Datagram.Counter {
				return nil
			}
		case <-time.After(delay):
			retries++
			delay *= 2 // Exponential backoff
			fmt.Printf("Timeout waiting for ACK, retrying... (%d/%d)\n", retries, maxRetries)
		}
	}

	session.ackRegistry.cleanupAck(ack)
	return fmt.Errorf("packet retransmission failed after %d attempts", maxRetries)
}

func handleAccountPeerComm(sm *SyncManager, session *Session, datagram Datagram) {
	// Create a key for the mutex based on username, peer username, and peer server address
	key := generateBaseKey(datagram.Username, datagram.PeerUsername, datagram.PeerServerAddress)
	mutex := sm.getMutex(key)
	mutex.Lock()
	defer mutex.Unlock()
	
	// Serialize the datagram to []byte
	packet := serializeDatagram(datagram)
	
	// Send the packet with retries
	if err := sendPacketWithRetry(session, packet, 5); err != nil {
		fmt.Printf("Failed to send packet for %s to %s: %v\n", datagram.Username, datagram.PeerUsername, err)
	}
}

func createSession(datagram Datagram, conn *net.UDPConn, sourceAddr *net.UDPAddr, ackRegistry *AckRegistry) *Session {
	// Creates a session, with the Datagram, connection, source address, and ack registry
	return &Session{
		Datagram:    datagram,
		conn:        conn,
		sourceAddr:  sourceAddr,
		ackRegistry: ackRegistry,
	}
}

func main() {
	// Listen on all interfaces, port 2012
	addr := net.UDPAddr{Port: 2012, IP: net.ParseIP("0.0.0.0")}
	conn, _ := net.ListenUDP("udp", &addr)
	defer conn.Close()

	ackRegistry := NewAckRegistry()
	syncManager := NewSyncManager()

	dispatcher := NewCentralDispatcher(conn, ackRegistry, syncManager)

	// Run the central dispatcher to listen for incoming packets
	go dispatcher.ListenAndServe()

	// Example usage of handleAccountPeerComm
	datagram := Datagram{
		Command:           0x02,
		Username:          "account1",
		PeerUsername:      "peer1",
		PeerServerAddress: "192.168.1.1",
		Counter:           1,
	}
	session := createSession(datagram, conn, nil, ackRegistry) // For client connections, sourceAddr would be set
	handleAccountPeerComm(syncManager, session, datagram)
}

func deserializeAck(data []byte) Ack {
	// Simplified deserialization logic
	return Ack{}
}

func deserializeDatagram(data []byte) Datagram {
	// Simplified deserialization logic
	return Datagram{}
}

func serializeDatagram(d Datagram) []byte {
	// Simplified serialization logic
	return []byte{}
}
