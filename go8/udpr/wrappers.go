package udpr

import (
	"ripple/config"
)

// Retry levels based on importance
const (
	LowImportance    = 5  // 5 retries for low importance messages
	HighImportance   = 12 // 12 retries for high importance messages
)

// Wrapper for SendWithRetry using a fixed port from config
func Send(data []byte, destinationAddr string) error {
	const maxRetries = 5
	return SendWithRetry(data, destinationAddr, config.Port, maxRetries)
}

// Wrapper for SendAck that takes a Conn struct
func Ack(idBytes []byte, c *Conn) error {
	return SendAck(idBytes, c.UDPConn, c.addr)
}
