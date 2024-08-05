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

    // Prepare success response
    responseData := make([]byte, 4) // Allocate 4 bytes for the trustline amount
    binary.BigEndian.PutUint32(responseData, uint32(trustlineAmount)) // Convert the trustline amount to bytes

    // Send the success response back to the client
    if err := handlers.SendSuccessResponse(ctx, responseData); err != nil {
        fmt.Printf("Error sending success response: %v\n", err) // Log the error
        return
    }

    fmt.Println("Inbound trustline amount sent successfully.")
}
