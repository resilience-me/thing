package client

import (
    "net"
    "os"
    "path/filepath"
    "resilience/main"
)

// Handles fetching the outbound trustline information
func GetTrustlineOut(ctx main.HandlerContext) {
    peerDir, err := main.GetPeerDir(ctx.Datagram, ctx.Datagram)
    if err != nil {
        // Send error response to the client
        _ = main.SendErrorResponse(ctx, "Error getting peer directory for outbound trustline.")
        return
    }

    trustlineOutPath := filepath.Join(peerDir, "trustline", "trustline_out.txt")
    trustlineAmount, err := os.ReadFile(trustlineOutPath)
    if err != nil {
        // Send error response to the client
        _ = main.SendErrorResponse(ctx, "Error reading outbound trustline file.")
        return
    }

    // Send success response with trustline amount
    _, err = ctx.Conn.WriteToUDP(trustlineAmount, ctx.Addr)
    if err != nil {
        // Send error response to the client
        _ = main.SendErrorResponse(ctx, "Error sending outbound trustline amount.")
        return
    }
}
