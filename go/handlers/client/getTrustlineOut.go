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

// GetTrustlineOut handles fetching the outbound trustline information
func GetTrustlineOut(ctx main.HandlerContext) {
    username := string(ctx.Datagram.XUsername[:])

    // Check if the account exists using the username from the datagram
    if err := main.CheckAccountExists(username); err != nil {
        fmt.Printf("Error getting account directory: %v\n", err) // Log detailed error
        _ = handlers.SendErrorResponse(ctx, "Failed to get account directory.") // Send simpler error message
        return
    }

    peerDir := main.GetPeerDir(ctx.Datagram)

    // Check if the peer directory exists
    if err := main.CheckPeerExists(peerDir); err != nil {
        fmt.Printf("Error getting peer directory: %v\n", err) // Log detailed error
        _ = handlers.SendErrorResponse(ctx, "Failed to get peer directory.") // Send simpler error message
        return
    }

    // Verify the client's signature
    if err := main.VerifyClientSignature(ctx.Datagram); err != nil {
        fmt.Printf("Signature verification failed: %v\n", err) // Log detailed error
        _ = handlers.SendErrorResponse(ctx, "Signature verification failed.") // Send simpler error message
        return
    }

    trustlineOutPath := filepath.Join(peerDir, "trustline", "trustline_out.txt")
    trustlineAmountStr, err := os.ReadFile(trustlineOutPath)
    if err != nil {
        fmt.Printf("Error reading trustline file: %v\n", err) // Log the error
        _ = handlers.SendErrorResponse(ctx, "Error reading outbound trustline file.")
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

    fmt.Println("Outbound trustline amount sent successfully.")
}
