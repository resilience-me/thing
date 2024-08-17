package udpr

import (
	"ripple/config"
)

// Wrapper for SendWithRetry using a fixed port from config
func Send(data []byte, destinationAddr string) error {
	const maxRetries = 5
	return SendWithRetry(data, destinationAddr, config.Port, maxRetries)
}

// Wrapper for SendAck that takes a Conn struct
func Ack(idBytes []byte, c *Conn) error {
	return SendAck(c.UDPConn, c.addr, idBytes)
}
