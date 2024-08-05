package client

import (
    "encoding/binary"
    "fmt"
    "net"
    "os"
    "path/filepath"
    "strconv"

    "resilience/main"
    "resilience/handlers"
)

// Handles fetching the inbound trustline information
func GetTrustlineIn(ctx main.HandlerContext) {
    // Validate the client request
    if err := handlers.ValidateClientRequest(ctx); err != nil {
        fmt.Printf("Validation failed: %v\n", err) // Log detailed error
        // Error response has already been sent in ValidateClientRequest
        return
    }

    // Get the peer directory for the trustline
    peerDir := main.GetPeerDir(ctx.Datagram)

    trustlineInPath := filepath.Join(peerDir, "trustline", "trustline_in.txt")
    trustlineAmountStr, err := os.ReadFile(trustlineInPath)
    if err != nil {
        fmt.Printf("Error reading inbound trustline file: %v\n", err) // Log the error
        _ = handlers.SendErrorResponse(ctx, "Error reading inbound trustline file.")
        return
    }

    // Convert the string to an integer
    trustlineAmount, err := strconv.ParseUint(string(trustlineAmountStr), 10, 32)
    if err != nil {
        fmt.Printf("Error converting trustline amount to integer: %v\n", err) // Log the error
        _ = handlers.SendErrorResponse(ctx, "Error converting trustline amount to integer.")
        return
    }

    // Prepare response datagram
    var responseDg main.ResponseDatagram
    responseDg.Result[0] = 0 // Set success code
    copy(responseDg.Nonce[:], ctx.Datagram.Signature[:]) // Use the original signature as the nonce

    // Store the trustline amount as bytes in the response
    binary.BigEndian.PutUint32(responseDg.Result[1:], uint32(trustlineAmount)) // Convert back to bytes
    if err := main.SignResponseDatagram(&responseDg, string(ctx.Datagram.XUsername[:])) ; err != nil {
        fmt.Printf("Failed to sign response datagram: %v\n", err) // Log the error
        _ = handlers.SendErrorResponse(ctx, "Failed to sign response datagram.")
        return
    }

    // Send the response back to the client
    _, err = ctx.Conn.WriteToUDP(responseDg[:], ctx.Addr)
    if err != nil {
        fmt.Printf("Error sending inbound trustline amount: %v\n", err) // Log the error
        return
    }

    fmt.Println("Inbound trustline amount sent successfully.")
}
