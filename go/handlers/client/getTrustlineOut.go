package client

import (
    "net"
    "os"
    "path/filepath"

    "resilience/main"
    "resilience/handlers"
)

// Handles fetching the outbound trustline information
func GetTrustlineOut(ctx main.HandlerContext) {
    peerDir, err := main.GetPeerDir(ctx.Datagram)
    if err != nil {
        _ = handlers.SendErrorResponse(ctx, "Error getting peer directory for outbound trustline.")
        return // Optionally handle or log the error
    }

    trustlineOutPath := filepath.Join(peerDir, "trustline", "trustline_out.txt")
    trustlineAmount, err := os.ReadFile(trustlineOutPath)
    if err != nil {
        _ = handlers.SendErrorResponse(ctx, "Error reading outbound trustline file.")
        return // Optionally handle or log the error
    }

    _, err = conn.WriteToUDP(trustlineAmount, addr)
    if err != nil {
        _ = handlers.SendErrorResponse(ctx, "Error sending outbound trustline amount.")
        return // Optionally handle or log the error
    }
}
