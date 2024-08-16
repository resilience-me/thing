package main

import (
	"fmt"
)

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
