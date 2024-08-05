package client

import (
    "fmt"
    "net"
    "os"
    "path/filepath"
    "strconv"

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
    trustlineAmountStr, err := os.ReadFile(trustlineOutPath)
    if err != nil {
        _ = handlers.SendErrorResponse(ctx, "Error reading outbound trustline file.")
        return
    }

    // Convert the string to an integer
    trustlineAmount, err := strconv.ParseUint(string(trustlineAmountStr), 10, 32)
    if err != nil {
        _ = handlers.SendErrorResponse(ctx, "Error converting trustline amount to integer.")
        return
    }

    // Prepare response datagram
    var responseDg main.ResponseDatagram
    responseDg.Result[0] = 0 // Set success code
    copy(responseDg.Nonce[:], ctx.Datagram.Signature[:]) // Use the original signature as the nonce

    // Store the trustline amount as bytes in the response
    binary.BigEndian.PutUint32(responseDg.Result[1:], uint32(trustlineAmount)) // Convert back to bytes
    if err := main.SignResponseDatagram(&responseDg, peerDir); err != nil {
        _ = handlers.SendErrorResponse(ctx, "Failed to sign response datagram.")
        return
    }

    // Send the response back to the client
    _, err = ctx.Conn.WriteToUDP(responseDg[:], ctx.Addr)
    if err != nil {
        _ = handlers.SendErrorResponse(ctx, "Error sending outbound trustline amount.")
        return
    }
}
