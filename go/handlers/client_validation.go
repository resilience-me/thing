package handlers

import (
    "fmt"
    "resilience/main"
)

// ValidateClientRequest checks if the client request is valid by verifying the account and peer directories and the client's signature.
func ValidateClientRequest(dg *main.Datagram) (string, error) {
    // Check if the account exists using the username from the datagram
    if err := main.CheckAccountExists(dg); err != nil {
        return "Failed to get account directory.", fmt.Errorf("failed to get account directory: %w", err)
    }

    // Check if the peer directory exists
    if err := main.CheckPeerExists(dg); err != nil {
        return "Failed to get peer directory.", fmt.Errorf("failed to get peer directory: %w", err)
    }

    // Verify the client's signature
    if err := main.VerifyClientSignature(dg); err != nil {
        return "Signature verification failed.", fmt.Errorf("signature verification failed: %w", err)
    }

    return "", nil // No error message, nil indicates a valid request
}
