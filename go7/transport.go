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

type AckEntry struct {
	peerAccount string
	ch          chan struct{}
}

// AckRegistry manages ACKs for different accounts
type AckRegistry struct {
	mu          sync.Mutex
	waitingAcks map[string]*AckEntry
}

// NewAckRegistry creates a new AckRegistry
func NewAckRegistry() *AckRegistry {
	return &AckRegistry{
		waitingAcks: make(map[string]*AckEntry),
	}
}

type SendContext struct {
	Data            []byte
	DestinationAddr string
	AckKey          string
	AckRegistry     *AckRegistry
	MaxRetries      int
}

// RegisterAck registers an Ack and ensures only one active peerAccount channel per username
func (ar *AckRegistry) RegisterAck(username string, peerAccount string) chan struct{} {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	ch := make(chan struct{})
	ar.waitingAcks[username] = &AckEntry{peerAccount: peerAccount, ch: ch}
	return ch
}

// RouteAck routes an incoming ACK to the appropriate channel
func (ar *AckRegistry) RouteAck(username string, peerAccount string) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	if entry, exists := ar.waitingAcks[username]; exists {
		if entry.peerAccount == peerAccount {
			close(entry.ch) // Signal the receipt of the ACK by closing the channel
			delete(ar.waitingAcks, username)
		}
	}
}

// CleanupAck removes an ACK from the registry after processing or timeout
func (ar *AckRegistry) CleanupAck(username string) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	delete(ar.waitingAcks, username)
}

// SendWithRetry sends data with retransmission logic based on the provided SendContext
func SendWithRetry(ctx SendContext) error {
	retries := 0
	delay := 1 * time.Second

	// Create a new UDP connection for sending the datagram
	sendConn, err := net.DialUDP("udp", nil, ctx.DestinationAddr)
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %w", err)
	}
	defer sendConn.Close()

	// Register the ACK using the provided AckKey and PeerAccount
	ackChan := ctx.AckRegistry.RegisterAck(ctx.AckKey, ctx.PeerAccount)

	for retries < ctx.MaxRetries {
		// Send the serialized datagram
		if _, err := sendConn.Write(ctx.Data); err != nil {
			return fmt.Errorf("failed to send data to server '%s': %w", ctx.DestinationAddr.String(), err)
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

// send is a lower-level function that sends data to a specified address
func send(data []byte, destinationAddr string) error {
	// Resolve the destination address using the provided address
	addr, err := net.ResolveUDPAddr("udp", destinationAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve server address '%s': %w", destinationAddr, err)
	}

	// Create a new UDP connection for sending the data
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %w", err)
	}
	defer conn.Close()

	// Send the data
	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("failed to send data to server '%s': %w", destinationAddr, err)
	}

	return nil
}

// SendAck is a wrapper around the lower-level send function, specifically for ACKs
func SendAck(data []byte, destinationAddr string) error {
	return send(data, destinationAddr)
}
