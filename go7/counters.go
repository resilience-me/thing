package main

import (
	"fmt"
)

// ValidateCounter checks if the datagram's counter is valid by comparing it to the last known counter.
// Returns an error if the counter is invalid (replay attack or stale datagram).
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

	return nil // Counter is valid
}
