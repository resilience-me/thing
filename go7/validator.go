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

// validateClientDatagram validates the client datagram
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

    return "", nil
}

// validateServerDatagram validates the server datagram
func validateServerDatagram(buf []byte, dg *Datagram) error {
    secretKey, err := loadServerSecretKey(dg)
    if err != nil {
        return fmt.Errorf("loading server secret key failed: %w", err)
    }

    if !verifyHMAC(buf, secretKey) {
        return errors.New("HMAC verification failed")
    }

    return nil
}

// ValidateDatagram parses and validates the received datagram based on its type (client or server)
// It also validates the counter to prevent replay attacks.
func ValidateDatagram(buffer []byte) (*Datagram, error) {
	// Parse the datagram
	datagram := parseDatagram(buffer)

	// Validate the datagram based on its type (client or server)
	if datagram.Command&0x80 == 0 { // Server session if MSB is 0
		if err := validateServerDatagram(buffer, datagram); err != nil {
			return nil, fmt.Errorf("error validating server datagram: %v", err)
		}
	} else { // Client session if MSB is 1
		errorMessage, err := validateClientDatagram(buffer, datagram)
		if err != nil {
			return nil, fmt.Errorf("error during client datagram validation: %v", err)
		}
	}

	// Validate the counter to prevent replay attacks
	if err := ValidateCounter(datagram); err != nil {
		return nil, fmt.Errorf("invalid counter detected: %v", err)
	}

	// Return the validated datagram
	return datagram, nil
}
