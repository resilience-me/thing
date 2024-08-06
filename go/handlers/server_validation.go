package handlers

import (
    "fmt"
    "resilience/main"
)

// ValidateServerRequest validates the server request by checking account existence, peer existence, and signature verification.
func ValidateServerRequest(dg *main.Datagram) error {
    // Check if the account exists
    if err := main.CheckAccountExists(dg); err != nil {
        return fmt.Errorf("error checking account existence: %v", err)
    }

    // Check if the peer exists
    if err := main.CheckPeerExists(dg); err != nil {
        return fmt.Errorf("error checking peer existence: %v", err)
    }

    // Verify the signature
    if err := main.VerifyServerSignature(dg); err != nil {
        return fmt.Errorf("signature verification failed: %v", err)
    }

    return nil // nil indicates a valid request
}
