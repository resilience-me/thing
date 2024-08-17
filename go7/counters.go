package main

import (
	"fmt"
)

// ValidateCounter checks if the datagram's counter is valid by comparing it to the last known counter.
// If valid, it sets the counter to the value in the datagram to prevent replay attacks.
func ValidateCounter(datagram *Datagram) error {
	var lastCounter uint32
	var counterError error

	// Determine the correct counter based on the command type
	if datagram.Command == 0x00 { // ACK datagram
		lastCounter, counterError = GetCounterOut(datagram)
	} else if datagram.Command&0x80 == 0 { // MSB is 0: Server connection
		lastCounter, counterError = GetCounterIn(datagram)
	} else { // MSB is 1: Client connection
		lastCounter, counterError = GetCounter(datagram)
	}

	if counterError != nil {
		return fmt.Errorf("error retrieving counter: %v", counterError)
	}

	// Validate the counter to prevent replay attacks
	if datagram.Counter <= lastCounter {
		return fmt.Errorf("replay detected or old datagram: Counter %d is not greater than the last seen counter %d", datagram.Counter, lastCounter)
	}

	// Set the counter to the value in the datagram to prevent replay attacks
	if datagram.Command == 0x00 { // ACK datagram
	} else if datagram.Command&0x80 == 0 { // Server connection
		err := SetCounterIn(datagram)
		if err != nil {
			return fmt.Errorf("failed to set in counter: %v", err)
		}
	} else { // Client connection
		err := SetCounter(datagram)
		if err != nil {
			return fmt.Errorf("failed to set counter: %v", err)
		}
	}

	return nil // Counter is valid and set
}
