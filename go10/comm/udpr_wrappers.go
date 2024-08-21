package comm

import (
	"ripple/config"
	"ripple/udpr"
)

// Retry levels based on importance
const (
	LowImportance    = 5  // 5 retries for standard messages
	HighImportance   = 12 // 12 retries for priority messages
)

// SendWithResolvedAddressAndConn resolves the address, creates a new UDP connection, and sends data with retries.
func SendWithResolvedAddressAndConn(address string, data []byte, maxRetries int) error {
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
	return udpr.SendWithRetry(conn, addr, data, maxRetries)
}

// Default Send with standard importance (5 retries)
func Send(destinationAddr string, data []byte) error {
	return SendWithResolvedAddressAndConn(destinationAddr, data, LowImportance)
}

// Send with priority importance (12 retries)
func SendPriority(destinationAddr string, data []byte) error {
	return SendWithResolvedAddressAndConn(destinationAddr, data, HighImportance)
}

// Wrapper for udpr.SendAck
func SendAck(conn *net.UDPConn, addr *net.UDPAddr, idBytes []byte) error {
	return udpr.SendAck(conn, addr, idBytes)
}
