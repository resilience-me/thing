package main

import (
	"fmt"
	"net"
	"time"
)

// SendWithRetry sends data with retransmission logic and continuously listens for an acknowledgment
func SendWithRetry(data []byte, destinationAddr string, maxRetries int) error {
	ackChan := make(chan bool, 1)

	// Resolve the destination address to a UDP address
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", destinationAddr, Port))
	if err != nil {
		return fmt.Errorf("failed to resolve server address '%s': %w", destinationAddr, err)
	}

	// Create a new UDP connection for sending the datagram
	sendConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %w", err)
	}
	defer sendConn.Close()

	// Start a goroutine to listen for the acknowledgment
	go listenForAck(sendConn, ackChan)

	retries := 0
	delay := 1 * time.Second

	for retries < maxRetries {
		// Send the datagram
		if _, err := sendConn.Write(data); err != nil {
			return fmt.Errorf("failed to send data to server '%s': %w", addr.String(), err)
		}

		select {
		case ackReceived := <-ackChan:
			if ackReceived {
				return nil // ACK received successfully, exit the function
			}
		case <-time.After(delay):
			retries++
			delay *= 2 // Exponential backoff
			fmt.Printf("Timeout or invalid ACK, retrying... (%d/%d)\n", retries, maxRetries)
		}
	}

	return fmt.Errorf("retransmission failed after %d attempts", maxRetries)
}

// listenForAck continuously listens for an acknowledgment
func listenForAck(conn *net.UDPConn, ackChan chan bool) {
	for {
		ack := make([]byte, 1)
		_, _, err := conn.ReadFromUDP(ack)
		if err == nil && ack[0] == AckByte {
			ackChan <- true
			return
		}
	}
}

// SendAck sends a simple acknowledgment (0xFF) using the provided Conn object
func SendAck(conn *Conn) error {
	ack := []byte{AckCode} // ACK value

	if _, err := conn.conn.WriteToUDP(ack, conn.addr); err != nil {
		return fmt.Errorf("failed to send ACK: %w", err)
	}
	return nil
}
