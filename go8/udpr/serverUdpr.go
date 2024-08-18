package udpr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

// SendWithRetryServer sends data with retries and waits for an acknowledgment using direct check
func SendWithRetryServer(conn *net.UDPConn, addr *net.UDPAddr, data []byte, maxRetries int) error {
	idBytes := newAck()
	return sendWithRetry(conn, addr, data, idBytes, maxRetries, func(delay time.Duration) bool {
		ack := make([]byte, 4)
		conn.SetReadDeadline(time.Now().Add(delay)) // Set the timeout for the read operation
		_, _, err := conn.ReadFromUDP(ack)
		if err != nil {
			return false
		}
		return bytes.Equal(ack, idBytes)
	})
}

// SendAck sends a simple acknowledgment with the byte slice identifier
func SendAck(conn *net.UDPConn, addr *net.UDPAddr, idBytes []byte) error {
	// Directly send the identifier as the ACK
	if _, err := conn.WriteToUDP(idBytes, addr); err != nil {
		return fmt.Errorf("failed to send ACK: %w", err)
	}
	return nil
}
