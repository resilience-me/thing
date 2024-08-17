package udpr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

// Global counter for generating unique 32-bit identifiers
var identifierCounter uint32

// SendWithRetry sends data with retransmission logic and waits for a simple acknowledgment
func SendWithRetry(data []byte, destinationAddr string, port int, maxRetries int) error {
	retries := 0
	delay := 1 * time.Second

	// Generate a unique 32-bit identifier for this transmission
	identifier := atomic.AddUint32(&identifierCounter, 1)

	// Convert the identifier to a 4-byte slice
	idBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(idBytes, identifier)

	// Resolve the destination address to a UDP address
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", destinationAddr, port))
	if err != nil {
		return fmt.Errorf("failed to resolve server address '%s': %w", destinationAddr, err)
	}

	// Create a single UDP connection for all attempts
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %w", err)
	}
	defer conn.Close()

	// Create the packet with the 4-byte identifier
	packet := append(idBytes, data...)

	for retries < maxRetries {
		// Send the datagram with the identifier
		if _, err := conn.Write(packet); err != nil {
			return fmt.Errorf("failed to send data to server '%s': %w", addr.String(), err)
		}

		// Set a deadline for the read operation
		conn.SetReadDeadline(time.Now().Add(delay))

		// Wait for the acknowledgment
		ack := make([]byte, 4)
		_, _, err = conn.ReadFromUDP(ack)

		if err == nil && bytes.Equal(ack, idBytes) {
			// Correct ACK received successfully
			return nil
		}

		// No correct ACK or an error occurred, retry
		retries++
		delay *= 2 // Exponential backoff
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
