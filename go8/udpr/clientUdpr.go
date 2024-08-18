package udpr

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
	"sync/atomic"
)

// AckManager manages acknowledgment registrations and checks
type AckManager struct {
	mu           sync.Mutex
	ackRegistry  map[string]struct{}
}

// NewAckManager initializes a new AckManager
func NewAckManager() *AckManager {
	return &AckManager{
		ackRegistry: make(map[string]struct{}),
	}
}

// SendWithRetryClient sends data with retries and waits for an acknowledgment using AckManager
func SendWithRetryClient(c *Client, data []byte, maxRetries int) error {
	idBytes := newAck()
	registerAck(c.ackManager, idBytes)

	return sendWithRetry(c.UDPConn, c.addr, packet, idBytes, maxRetries, func(delay time.Duration) bool {
		// Wait for an ACK or timeout within the checkAck function
		time.Sleep(delay)

		// Check if the ACK has been received
		return pollAck(c.ackManager, idBytes)
	})
}

// SendWithRetryClient sends data with retries and waits for an acknowledgment
func SendWithRetryClient(c *Client, data []byte, maxRetries int) error {
	delay := 1 * time.Second
	maxDelay := 10 * time.Second

	identifier := atomic.AddUint32(&identifierCounter, 1)
	idBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(idBytes, identifier)

	// Register the identifier in the ackRegistry
	registerAck(c.ackManager, idBytes)

	packet := append(idBytes, data...)
	for retries := 0; retries < maxRetries; retries++ {
		_, err := c.UDPConn.WriteToUDP(packet, c.addr)
		if err != nil {
			return fmt.Errorf("failed to send data: %w", err)
		}

		// Wait for an ACK or timeout
		time.Sleep(delay)

		// Check if the ACK has been received
		if pollAck(c.ackManager, idBytes) {
			// ACK received
			return nil
		}

		// ACK not received, retry
		fmt.Println("Retrying...")
		retries++
		if delay < maxDelay {
			delay *= 2 // Exponential backoff
		}
		fmt.Printf("Timeout or invalid ACK, retrying... (%d/%d)\n", retries, maxRetries)
	}

	return fmt.Errorf("retransmission failed after %d attempts", maxRetries)
}

// registerAck adds the ACK identifier to the registry
func registerAck(ackMgr *AckManager, idBytes []byte) {
	ackMgr.mu.Lock()
	ackMgr.ackRegistry[string(idBytes)] = struct{}{}
	ackMgr.mu.Unlock()
}

// pollAck checks if the ACK identifier is present in the registry
func pollAck(ackMgr *AckManager, idBytes []byte) bool {
	ackMgr.mu.Lock()
	_, exists := ackMgr.ackRegistry[string(idBytes)]
	ackMgr.mu.Unlock()
	return !exists // Return true if ACK is NOT present (i.e., we should retry)
}
