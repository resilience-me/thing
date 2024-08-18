package udpr

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
	"sync/atomic"
)

func SendWithRetryClient(c *Client, data []byte, maxRetries int) error {
	delay := 1 * time.Second
	maxDelay := 10 * time.Second

	identifier := atomic.AddUint32(&identifierCounter, 1)
	idBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(idBytes, identifier)

	// Register the identifier in the ackRegistry
	c.ackManager.mu.Lock()
	c.ackManager.ackRegistry[string(idBytes)] = struct{}{}
	c.ackManager.mu.Unlock()

	packet := append(idBytes, data...)
	for retries := 0; retries < maxRetries; retries++ {
		_, err := c.UDPConn.WriteToUDP(packet, c.addr)
		if err != nil {
			return fmt.Errorf("failed to send data: %w", err)
		}

		// Wait for an ACK or timeout
		time.Sleep(delay)

		// Check if the ACK has been received
		c.ackManager.mu.Lock()
		exists := c.ackManager.ackRegistry[string(idBytes)]
		c.ackManager.mu.Unlock()

		if !exists {
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
