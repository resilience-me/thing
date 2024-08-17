package main

import (
	"fmt"
	"net"
	"time"
)

// SendWithRetry sends data with retransmission logic and waits for a simple acknowledgment
func SendWithRetry(data []byte, destinationAddr string, maxRetries int) error {
	retries := 0
	delay := 1 * time.Second

	// Resolve the destination address to a UDP address
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", destinationAddr, Port))
	if err != nil {
		return fmt.Errorf("failed to resolve server address '%s': %w", destinationAddr, err)
	}

	// Create a single UDP connection for all attempts
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %w", err)
	}
	defer conn.Close()

	for retries < maxRetries {
		// Send the datagram
		if _, err := conn.Write(data); err != nil {
			return fmt.Errorf("failed to send data to server '%s': %w", addr.String(), err)
		}

		// Set a deadline for the read operation
		conn.SetReadDeadline(time.Now().Add(delay))

		// Wait for the acknowledgment
		ack := make([]byte, 1)
		_, _, err = conn.ReadFromUDP(ack)

		if err == nil && ack[0] == AckCode {
			// ACK received successfully
			return nil
		}

		// No ACK or an error occurred, retry
		retries++
		delay *= 2 // Exponential backoff
		fmt.Printf("Timeout or invalid ACK, retrying... (%d/%d)\n", retries, maxRetries)
	}

	return fmt.Errorf("retransmission failed after %d attempts", maxRetries)
}

// SendAck sends a simple acknowledgment (0xFF) using the provided Conn object
func SendAck(conn *Conn) error {
	ack := []byte{AckCode} // ACK value

	if _, err := conn.WriteToUDP(ack, conn.addr); err != nil {
		return fmt.Errorf("failed to send ACK: %w", err)
	}
	return nil
}
