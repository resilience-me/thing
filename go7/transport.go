package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Ack represents an acknowledgment packet
type Ack struct {
	Username          string
	PeerUsername      string
	PeerServerAddress string
	Counter           uint32
}

// NewAck creates a new Ack from a given Datagram
func NewAck(dg *Datagram) *Ack {
	return &Ack{
		Username:          dg.Username,
		PeerUsername:      dg.PeerUsername,
		PeerServerAddress: dg.PeerServerAddress,
		Counter:           dg.Counter,
	}
}

// AckRegistry manages ACKs for different accounts
type AckRegistry struct {
	mu          sync.Mutex
	waitingAcks map[string]chan *Ack
}

// NewAckRegistry creates a new AckRegistry
func NewAckRegistry() *AckRegistry {
	return &AckRegistry{
		waitingAcks: make(map[string]chan *Ack),
	}
}

// RegisterAck registers an Ack and returns a channel to receive it
func (ar *AckRegistry) RegisterAck(ack *Ack) chan *Ack {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	key := generateAckKey(ack)
	ch := make(chan *Ack)
	ar.waitingAcks[key] = ch
	return ch
}

// RouteAck routes an incoming ACK to the appropriate channel
func (ar *AckRegistry) RouteAck(ack *Ack) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	key := generateAckKey(ack)
	if ch, exists := ar.waitingAcks[key]; exists {
		ch <- ack
		close(ch) // Close the channel after sending the ACK to avoid leaks
		delete(ar.waitingAcks, key)
	}
}

// CleanupAck removes an ACK from the registry after processing or timeout
func (ar *AckRegistry) CleanupAck(ack *Ack) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	key := generateAckKey(ack)
	delete(ar.waitingAcks, key)
}

// Transport provides reliable transmission functionality over UDP
type Transport struct {
	conn        *net.UDPConn
	ackRegistry *AckRegistry
}

// NewTransport creates a new Transport instance
func NewTransport(conn *net.UDPConn, ackRegistry *AckRegistry) *Transport {
	return &Transport{
		conn:        conn,
		ackRegistry: ackRegistry,
	}
}

// SendPacketWithRetry sends a packet with retransmission logic
func (t *Transport) SendPacketWithRetry(session *Session, packet []byte, maxRetries int) error {
	retries := 0
	delay := 1 * time.Second

	serverAddress := session.Datagram.PeerServerAddress
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", serverAddress, 2012))
	if err != nil {
		return fmt.Errorf("failed to resolve server address '%s': %w", serverAddress, err)
	}

	ack := NewAck(session.Datagram)

	ackChan := t.ackRegistry.RegisterAck(ack)

	for retries < maxRetries {
		if _, err := t.conn.WriteToUDP(packet, addr); err != nil {
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

	t.ackRegistry.CleanupAck(ack)
	return fmt.Errorf("packet retransmission failed after %d attempts", maxRetries)
}

// Utility function to generate a unique key for ACKs based on the Ack fields
func generateAckKey(ack *Ack) string {
	return fmt.Sprintf("%s-%s-%s-%d", ack.Username, ack.PeerUsername, ack.PeerServerAddress, ack.Counter)
}
