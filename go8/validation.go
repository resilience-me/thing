package main

import (
	"fmt"
)

// checkUserAndPeerExist checks for the existence of user and peer directories
// It returns an error message string (empty if successful) and an error object for detailed information if an error occurs.
func checkUserAndPeerExist(dg *Datagram) (string, error) {
    exists, err := database.CheckAccountExists(dg)
    if err != nil {
        return "Error checking account existence", fmt.Errorf("error checking account existence for user '%s': %v", dg.Username, err)
    } else if !exists {
        return "User account does not exist", fmt.Errorf("account directory does not exist for user '%s'", dg.Username)
    }

    exists, err = database.CheckPeerExists(dg)
    if err != nil {
        return "Error checking peer existence", fmt.Errorf("error checking peer existence for server '%s' and user '%s': %v", dg.PeerServerAddress, dg.PeerUsername, err)
    } else if !exists {
        return "Peer account does not exist", fmt.Errorf("peer directory does not exist for server '%s' and user '%s'", dg.PeerServerAddress, dg.PeerUsername)
    }

    return "", nil // No error, directories exist
}

// validateClientDatagram validates the client datagram and checks the counter
func validateClientDatagram(buf []byte, dg *Datagram) (string, error) {
    errorMessage, err := checkUserAndPeerExist(dg)
    if err != nil {
        return errorMessage, fmt.Errorf("user and peer existence check failed: %w", err)
    }

    secretKey, err := loadClientSecretKey(dg)
    if err != nil {
        return "Error loading client secret key", fmt.Errorf("loading client secret key failed: %w", err)
    }

    if !verifyHMAC(buf, secretKey) {
        return "Error verifying HMAC", errors.New("HMAC verification failed")
    }

    // Validate the counter
    if err := ValidateAndIncrementClientCounter(dg); err != nil {
        return "Invalid counter", fmt.Errorf("counter validation failed: %w", err)
    }

    return "", nil
}

// validateServerDatagram validates the server datagram and checks the counter
func validateServerDatagram(buf []byte, dg *Datagram) error {
    secretKey, err := loadServerSecretKey(dg)
    if err != nil {
        return fmt.Errorf("loading server secret key failed: %w", err)
    }

    if !verifyHMAC(buf, secretKey) {
        return errors.New("HMAC verification failed")
    }

    // Validate the counter
    if err := ValidateAndIncrementServerCounter(dg); err != nil {
        return fmt.Errorf("counter validation failed: %w", err)
    }

    return nil
}
