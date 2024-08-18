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
