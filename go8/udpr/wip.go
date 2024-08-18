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
func sendWithRetry(conn *net.UDPConn, addr *net.UDPAddr, packet []byte, maxRetries int, checkAck func(delay time.Duration) bool) error {
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
func newAck() []byte {
	identifier := atomic.AddUint32(&identifierCounter, 1)
	idBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(idBytes, identifier)
	return idBytes
}
