package main

import (
	"fmt"
)

// checkPeerExist checks for the existence of user and peer directories
// It returns an error message string (empty if successful) and an error object for detailed information if an error occurs.
func checkPeerExist(dg *Datagram) (string, error) {
    exists, err = database.CheckPeerExists(dg)
    if err != nil {
        return "Error checking peer existence", fmt.Errorf("error checking peer existence for server '%s' and user '%s': %v", dg.PeerServerAddress, dg.PeerUsername, err)
    } else if !exists {
        return "Peer account does not exist", fmt.Errorf("peer directory does not exist for server '%s' and user '%s'", dg.PeerServerAddress, dg.PeerUsername)
    }

    return "", nil // No error, directories exist
}

// validateClientDatagram validates the client datagram and checks the counter
func validateClientDatagram(buf []byte, dg *Datagram) error {

    secretKey, err := loadClientSecretKey(dg)
    if err != nil {
        return fmt.Errorf("loading client secret key failed: %w", err)
    }

    if !verifyHMAC(buf, secretKey) {
        return errors.New("HMAC verification failed")
    }

    // Validate the counter
    if err := ValidateAndIncrementClientCounter(dg); err != nil {
        return fmt.Errorf("counter validation failed: %w", err)
    }

    return nil
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

// validateDatagram validates a datagram based on whether it's for a server or client session.
func validateDatagram(buf []byte, dg *Datagram) error {
	if dg.Command&0x80 == 0 { // Server session if MSB is 0
		return validateServerDatagram(buf, dg)
	} else {  // Client session if MSB is 1
		return validateClientDatagram(buf, dg)
	}
}
