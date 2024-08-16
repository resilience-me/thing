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

type SendContext struct {
	Data           []byte
	DestinationAddr string
	AckKey         string
	AckRegistry    *AckRegistry
	MaxRetries     int
}

// RegisterAck registers an Ack and returns a channel to receive the trigger signal
func (ar *AckRegistry) RegisterAck(ackKey string) chan struct{} {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	ch := make(chan struct{})
	ar.waitingAcks[ackKey] = ch
	return ch
}

// RouteAck routes an incoming ACK to the appropriate channel
func (ar *AckRegistry) RouteAck(ackKey string) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	if ch, exists := ar.waitingAcks[ackKey]; exists {
		close(ch) // Signal the receipt of the ACK by closing the channel
		delete(ar.waitingAcks, ackKey)
	}
}

// CleanupAck removes an ACK from the registry after processing or timeout
func (ar *AckRegistry) CleanupAck(ackKey string) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	delete(ar.waitingAcks, ackKey)
}

// SendWithRetry sends data with retransmission logic based on the provided SendContext
func SendWithRetry(ctx SendContext) error {
	retries := 0
	delay := 1 * time.Second

	// Resolve the destination address to a UDP address
	addr, err := net.ResolveUDPAddr("udp", ctx.DestinationAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve server address '%s': %w", ctx.DestinationAddr, err)
	}

	// Create a new UDP connection for sending the datagram
	sendConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %w", err)
	}
	defer sendConn.Close()

	// Register the ACK using the provided AckKey
	ackChan := ctx.AckRegistry.RegisterAck(ctx.AckKey)

	for retries < ctx.MaxRetries {
		// Send the serialized datagram
		if _, err := sendConn.Write(ctx.Data); err != nil {
			return fmt.Errorf("failed to send data to server '%s': %w", ctx.DestinationAddr, err)
		}

		select {
		case <-ackChan: // Waiting for the channel to be closed as a signal
			return nil
		case <-time.After(delay):
			retries++
			delay *= 2 // Exponential backoff
			fmt.Printf("Timeout waiting for ACK, retrying... (%d/%d)\n", retries, ctx.MaxRetries)
		}
	}

	// Cleanup the ACK registration if we failed to get the ACK
	ctx.AckRegistry.CleanupAck(ctx.AckKey)
	return fmt.Errorf("retransmission failed after %d attempts", ctx.MaxRetries)
}

// Utility function to generate a unique key for ACKs based on the Ack fields
func generateAckKey(ack *Ack) string {
	return fmt.Sprintf("%s-%s-%s-%d", ack.Username, ack.PeerUsername, ack.PeerServerAddress, ack.Counter)
}
