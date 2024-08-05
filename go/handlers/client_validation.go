package handlers

import (
    "fmt"
    "resilience/main"
)

// ValidateClientRequest checks if the client request is valid by verifying the account and peer directories and the client's signature.
func ValidateClientRequest(ctx main.HandlerContext) error {
    // Check if the account exists using the username from the datagram
    if err := main.CheckAccountExists(ctx.Datagram); err != nil {
        fmt.Printf("Error getting account directory: %v\n", err) // Log detailed error
        _ = SendErrorResponse(ctx, "Failed to get account directory.") // Send simpler error message
        return err
    }

    // Check if the peer directory exists
    if err := main.CheckPeerExists(ctx.Datagram); err != nil {
        fmt.Printf("Error getting peer directory: %v\n", err) // Log detailed error
        _ = SendErrorResponse(ctx, "Failed to get peer directory.") // Send simpler error message
        return err
    }

    // Verify the client's signature
    if err := main.VerifyClientSignature(ctx.Datagram); err != nil {
        fmt.Printf("Signature verification failed: %v\n", err) // Log detailed error
        _ = SendErrorResponse(ctx, "Signature verification failed.") // Send simpler error message
        return err
    }

    return nil // nil indicates a valid request
}
