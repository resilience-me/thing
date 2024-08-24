package comm

import (
	"fmt"
	"net"
	"ripple/config"
	"ripple/udpr"
)

// Retry levels based on importance
const (
	LowImportance    = 5  // 5 retries for standard messages
	HighImportance   = 12 // 12 retries for priority messages
)

// SendWithAddress sends data to a specified UDP address with retry logic.
// It handles the creation and closure of the UDP connection internally.
func SendWithAddress(addr *net.UDPAddr, data []byte, maxRetries int) error {
	// Create a UDP connection with an ephemeral local port
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %w", err)
	}
	defer conn.Close()

	// Use the udpr.SendWithRetry to send the data
	if err := udpr.SendWithRetry(conn, addr, data, maxRetries); err != nil {
		return fmt.Errorf("error sending data: %w", err)
	}

	return nil
}

// SendWithResolvedAddressAndConn resolves the address, creates a new UDP connection, and sends data with retries.
func SendWithResolvedAddress(address string, data []byte, maxRetries int) error {
	// Resolve the destination address to a UDP address
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, config.Port))
	if err != nil {
		return fmt.Errorf("failed to resolve address '%s:%d': %w", address, config.Port, err)
	}
	// Call SendWithAddress function with the resolved address
	return SendWithAddress(addr, data, maxRetries)
}

// // Default Send with standard importance (5 retries)
// func Send(destinationAddr string, data []byte) error {
// 	return SendWithResolvedAddress(destinationAddr, data, LowImportance)
// }

// // Send with priority importance (12 retries)
// func SendPriority(destinationAddr string, data []byte) error {
// 	return SendWithResolvedAddress(destinationAddr, data, HighImportance)
// }

// Wrapper for udpr.SendAck
func SendAck(conn *net.UDPConn, addr *net.UDPAddr, idBytes []byte) error {
	return udpr.SendAck(conn, addr, idBytes)
}
