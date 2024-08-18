package udpr

import (
	"ripple/config"
)

// Retry levels based on importance
const (
	LowImportance    = 5  // 5 retries for standard messages
	HighImportance   = 12 // 12 retries for priority messages
)

// SendWithResolvedAddressAndConn resolves the address, creates a new UDP connection, and sends data with retries.
func SendWithResolvedAddressAndConn(address string, data []byte, maxRetries int) error {
	// Resolve the destination address to a UDP address
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return fmt.Errorf("failed to resolve address '%s:%d': %w", address, port, err)
	}

	// Create a UDP connection with an ephemeral local port
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %w", err)
	}
	defer conn.Close()

	// Call the SendWithRetry function with the resolved address and the newly created connection
	return SendWithRetry(conn, addr, data, maxRetries)
}

// Default Send with standard importance (5 retries)
func Send(data []byte, destinationAddr string) error {
	return SendWithRetry(data, destinationAddr, config.Port, LowImportance)
}

// Send with priority importance (12 retries)
func SendPriority(data []byte, destinationAddr string) error {
	return SendWithRetry(data, destinationAddr, config.Port, HighImportance)
}

// Wrapper for SendAck that takes a Conn struct
func Ack(idBytes []byte, c *Conn) error {
	return SendAck(idBytes, c.UDPConn, c.addr)
}
