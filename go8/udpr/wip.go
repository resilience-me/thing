package udpr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
	"sync/atomic"
)

// Client represents a UDP client with connection and ACK manager
type Client struct {
	UDPConn    *net.UDPConn
	addr       *net.UDPAddr
	ackManager *AckManager
}

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
	packet, idBytes := preparePacket(data)
	registerAck(c.ackManager, idBytes)

	return sendWithRetry(c.UDPConn, c.addr, packet, idBytes, maxRetries, func(delay time.Duration) bool {
		// Wait for an ACK or timeout within the checkAck function
		time.Sleep(delay)

		// Check if the ACK has been received
		return pollAck(c.ackManager, idBytes)
	})
}

// SendWithRetry sends data with retries and waits for an acknowledgment using direct check
func SendWithRetry(conn *net.UDPConn, addr *net.UDPAddr, data []byte, maxRetries int) error {
	packet, idBytes := preparePacket(data)
	return sendWithRetry(conn, addr, packet, idBytes, maxRetries, func(delay time.Duration) bool {
		ack := make([]byte, 4)
		conn.SetReadDeadline(time.Now().Add(delay)) // Set the timeout for the read operation
		_, _, err := conn.ReadFromUDP(ack)
		if err != nil {
			return false
		}
		return bytes.Equal(ack, idBytes)
	})
}

const (
	initialDelay = 1 * time.Second	   // Initial delay duration
	maxDelay = 16 * time.Second 	   // Maximum delay duration
)

// Global counter for generating unique 32-bit identifiers
var identifierCounter uint32

// sendWithRetry sends data with retries and checks for acknowledgment using the provided check function
func sendWithRetry(conn *net.UDPConn, addr *net.UDPAddr, packet []byte, idBytes []byte, maxRetries int, checkAck func(delay time.Duration) bool) error {
	delay := initialDelay

	for retries := 0; retries < maxRetries; retries++ {
		_, err := conn.WriteToUDP(packet, addr)
		if err != nil {
			return fmt.Errorf("failed to send data: %w", err)
		}

		// Check for ACK with the provided check function
		if checkAck(delay) {
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

// preparePacket generates a unique identifier and prepares the packet
func preparePacket(data []byte) ([]byte, []byte) {
	identifier := atomic.AddUint32(&identifierCounter, 1)
	idBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(idBytes, identifier)
	packet := append(idBytes, data...)
	return packet, idBytes
}

// registerAck adds the ACK identifier to the registry
func registerAck(ackMgr *AckManager, idBytes []byte) {
	ackMgr.mu.Lock()
	defer ackMgr.mu.Unlock()
	ackMgr.ackRegistry[string(idBytes)] = struct{}{}
}

// pollAck checks if the ACK identifier is present in the registry
func pollAck(ackMgr *AckManager, idBytes []byte) bool {
	ackMgr.mu.Lock()
	defer ackMgr.mu.Unlock()
	_, exists := ackMgr.ackRegistry[string(idBytes)]
	return !exists // Return true if ACK is NOT present (i.e., we should retry)
}
