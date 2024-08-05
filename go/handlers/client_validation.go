package handlers

import (
    "fmt"
    "resilience/main"
)

// ValidateClientRequest checks if the client request is valid by verifying the account and peer directories and the client's signature.
func ValidateClientRequest(ctx main.HandlerContext) error {
    // Check if the account exists using the username from the datagram
    if err := main.CheckAccountExists(ctx.Datagram); err != nil {
        _ = SendErrorResponse(ctx, "Failed to get account directory.") // Send simpler error message
        return fmt.Errorf("failed to get account directory: %w", err)
    }

    // Check if the peer directory exists
    if err := main.CheckPeerExists(ctx.Datagram); err != nil {
        _ = SendErrorResponse(ctx, "Failed to get peer directory.") // Send simpler error message
        return fmt.Errorf("failed to get peer directory: %w", err)
    }

    // Verify the client's signature
    if err := main.VerifyClientSignature(ctx.Datagram); err != nil {
        _ = SendErrorResponse(ctx, "Signature verification failed.") // Send simpler error message
        return fmt.Errorf("signature verification failed: %w", err)
    }

    return nil // nil indicates a valid request
}
