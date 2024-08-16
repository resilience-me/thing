package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Transport provides reliable transmission functionality over UDP
type Transport struct {
	*AckRegistry
}

// NewTransport creates a new Transport instance with a new AckRegistry
func NewTransport() *Transport {
	return &Transport{AckRegistry: NewAckRegistry()}
}

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

// SendPacketWithRetry sends a packet with retransmission logic
func SendPacketWithRetry(session *Session, packet []byte, maxRetries int) error {
	retries := 0
	delay := 1 * time.Second

	serverAddress := session.Datagram.PeerServerAddress
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", serverAddress, 2012))
	if err != nil {
		return fmt.Errorf("failed to resolve server address '%s': %w", serverAddress, err)
	}

	// Create a new UDP connection for sending the packet
	sendConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %w", err)
	}
	defer sendConn.Close()

	ack := NewAck(session.Datagram)
	ackChan := session.AckRegistry.RegisterAck(ack)

	for retries < maxRetries {
		if _, err := sendConn.Write(packet); err != nil {
			return fmt.Errorf("failed to send data to server '%s': %w", serverAddress, err)
		}

		select {
		case <-ackChan:
			// ACK received, no need to compare counters as the registry ensures it's the correct one
			return nil
		case <-time.After(delay):
			retries++
			delay *= 2 // Exponential backoff
			fmt.Printf("Timeout waiting for ACK, retrying... (%d/%d)\n", retries, maxRetries)
		}
	}

	session.AckRegistry.CleanupAck(ack)
	return fmt.Errorf("packet retransmission failed after %d attempts", maxRetries)
}

// Utility function to generate a unique key for ACKs based on the Ack fields
func generateAckKey(ack *Ack) string {
	return fmt.Sprintf("%s-%s-%s-%d", ack.Username, ack.PeerUsername, ack.PeerServerAddress, ack.Counter)
}
