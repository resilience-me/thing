package handlers

import (
    "fmt"
    "resilience/main" // Correct import path based on the Go module name
)

// SignAndSendDatagram signs a Datagram and sends it via UDP.
func SignAndSendDatagram(ctx main.HandlerContext, dg *main.Datagram) error {
    // Sign the datagram
    if err := main.SignDatagram(dg); err != nil {
        return fmt.Errorf("failed to sign datagram: %v", err)
    }

    // Send the datagram back to the peer
    if _, err := ctx.Conn.WriteToUDP(dg[:], ctx.Addr); err != nil {
        return fmt.Errorf("error sending datagram: %v", err)
    }

    return nil
}
