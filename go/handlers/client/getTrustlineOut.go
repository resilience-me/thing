package client

import (
    "fmt"
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
        return
    }

    trustlineOutPath := filepath.Join(peerDir, "trustline", "trustline_out.txt")
    trustlineAmount, err := os.ReadFile(trustlineOutPath)
    if err != nil {
        _ = handlers.SendErrorResponse(ctx, "Error reading outbound trustline file.")
        return
    }

    // Prepare response datagram
    var responseDg main.ResponseDatagram
    responseDg.Result[0] = 0 // Set success code
    copy(responseDg.Nonce[:], ctx.Datagram.Signature[:]) // Use the original signature as the nonce
    copy(responseDg.Result[1:], trustlineAmount) // Copy the trustline amount directly

    if err := main.SignResponseDatagram(&responseDg, peerDir); err != nil {
        fmt.Printf("Failed to sign response datagram: %v\n", err) // Log detailed error
        _ = handlers.SendErrorResponse(ctx, "Failed to sign response datagram.")
        return
    }

    // Send the response back to the client
    _, err = ctx.Conn.WriteToUDP(responseDg[:], ctx.Addr)
    if err != nil {
        _ = handlers.SendErrorResponse(ctx, "Error sending outbound trustline amount.")
        return
    }
    
    fmt.Println("Sent outbound trustline amount to client.")
}
