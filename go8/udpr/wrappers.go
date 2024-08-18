package udpr

import (
	"fmt"
	"net"
	"ripple/config"
)

// Retry levels based on importance
const (
	LowImportance    = 5  // 5 retries for standard messages
	HighImportance   = 12 // 12 retries for priority messages
)

// SendWithResolvedAddress resolves the address, creates a new UDP connection, and sends data with retries.
func SendWithResolvedAddress(address string, data []byte, maxRetries int) error {
	// Resolve the destination address to a UDP address
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, config.Port))
	if err != nil {
		return fmt.Errorf("failed to resolve address '%s:%d': %w", address, config.Port, err)
	}

	// Create a UDP connection with an ephemeral local port
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %w", err)
	}
	defer conn.Close()

	// Call the SendWithRetry function with the resolved address and the newly created connection
	return SendWithRetryServer(conn, addr, data, maxRetries)
}

// Default Send with standard importance (5 retries)
func SendServer(destinationAddr string, data []byte) error {
	return SendWithResolvedAddress(destinationAddr, data, LowImportance)
}

// Send with priority importance (12 retries)
func SendPriorityServer(destinationAddr string, data []byte) error {
	return SendWithResolvedAddress(destinationAddr, data, HighImportance)
}

// Default Send with standard importance (5 retries)
func SendClient(client *Client, data []byte) error {
	return SendWithRetryClient(client, data, LowImportance)
}

// Send with priority importance (12 retries)
func SendPriorityClient(client *Client, data []byte) error {
	return SendWithRetryClient(client, data, HighImportance)
}

// Wrapper for SendAck
func Ack(conn *net.UDPConn, addr *net.UDPAddr, idBytes []byte) error {
	return SendAck(conn, addr, idBytes)
}
