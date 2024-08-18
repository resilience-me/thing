package udpr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

// SendWithRetryServer sends data with retries and waits for an acknowledgment using direct check
func SendWithRetryServer(conn *net.UDPConn, addr *net.UDPAddr, data []byte, maxRetries int) error {
	idBytes := newAck()
	return sendWithRetry(conn, addr, data, idBytes, maxRetries, func(delay time.Duration) bool {
		ack := make([]byte, 4)
		conn.SetReadDeadline(time.Now().Add(delay)) // Set the timeout for the read operation
		_, _, err := conn.ReadFromUDP(ack)
		if err != nil {
			return false
		}
		return bytes.Equal(ack, idBytes)
	})
}

// SendWithRetry sends data with retransmission logic and waits for a simple acknowledgment
func SendWithRetry(conn *net.UDPConn, addr *net.UDPAddr, data []byte, maxRetries int) error {
	retries := 0
	delay := 1 * time.Second

	// Generate a unique 32-bit identifier for this transmission
	identifier := atomic.AddUint32(&identifierCounter, 1)

	// Convert the identifier to a 4-byte slice
	idBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(idBytes, identifier)

	// Create the packet with the 4-byte identifier
	packet := append(idBytes, data...)

	for retries < maxRetries {
		// Send the datagram with the identifier
		if _, err := conn.WriteToUDP(packet, addr); err != nil {
			return fmt.Errorf("failed to send data to %s: %w", addr.String(), err)
		}

		// Set a deadline for the read operation
		conn.SetReadDeadline(time.Now().Add(delay))

		// Wait for the acknowledgment
		ack := make([]byte, 4)
		_, _, err := conn.ReadFromUDP(ack)

		if err == nil && bytes.Equal(ack, idBytes) {
			// Correct ACK received successfully
			return nil
		}

		// No correct ACK or an error occurred, retry
		retries++
		if delay < maxDelay {
			delay *= 2 // Exponential backoff
		}
		fmt.Printf("Timeout or invalid ACK, retrying... (%d/%d)\n", retries, maxRetries)
	}

	return fmt.Errorf("retransmission failed after %d attempts", maxRetries)
}

// SendAck sends a simple acknowledgment with the byte slice identifier
func SendAck(conn *net.UDPConn, addr *net.UDPAddr, idBytes []byte) error {
	// Directly send the identifier as the ACK
	if _, err := conn.WriteToUDP(idBytes, addr); err != nil {
		return fmt.Errorf("failed to send ACK: %w", err)
	}
	return nil
}
