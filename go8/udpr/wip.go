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

const (
	initialDelay = 1 * time.Second	   // Initial delay duration
	maxDelay = 16 * time.Second 	   // Maximum delay duration
)

// Global counter for generating unique 32-bit identifiers
var identifierCounter uint32

// sendWithRetry sends data with retries and checks for acknowledgment using the provided check function
func sendWithRetry(conn *net.UDPConn, addr *net.UDPAddr, data []byte, idBytes []byte, maxRetries int, checkAck func(delay time.Duration) bool) error {
	delay := initialDelay

	// Create the packet with the 4-byte identifier
	packet := append(idBytes, data...)

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
		if delay < maxDelay {
			delay *= 2 // Exponential backoff
		}
		fmt.Printf("Timeout or invalid ACK, retrying... (%d/%d)\n", retries, maxRetries)
	}

	return fmt.Errorf("retransmission failed after %d attempts", maxRetries)
}

// preparePacket generates a unique identifier and prepares the packet
func newAck() []byte {
	identifier := atomic.AddUint32(&identifierCounter, 1)
	idBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(idBytes, identifier)
	return idBytes
}
