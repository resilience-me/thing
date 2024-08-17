package udpr

import (
	"ripple/config"
)

// Retry levels based on importance
const (
	LowImportance    = 5  // 5 retries for standard messages
	HighImportance   = 12 // 12 retries for priority messages
)

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
