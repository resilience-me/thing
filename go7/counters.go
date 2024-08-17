package main

import (
	"fmt"
)

// ValidateClientCounter checks if the datagram's counter is valid by comparing it to the last known counter for client connections.
// If valid, it sets the counter to the value in the datagram to prevent replay attacks.
func ValidateClientCounter(datagram *Datagram) error {
	prevCounter, err := GetCounter(datagram)
	if err != nil {
		return fmt.Errorf("error retrieving counter: %v", err)
	}
	if datagram.Counter < prevCounter {
		return fmt.Errorf("replay detected or old datagram: Counter %d is not greater than the last seen counter %d", datagram.Counter, prevCounter)
	}

	return nil
}

// ValidateServerCounter checks if the datagram's counter is valid by comparing it to the last known counter for server connections.
// If valid, it sets the counter to the value in the datagram to prevent replay attacks.
func ValidateServerCounter(datagram *Datagram) error {
	prevCounter, err := GetCounterIn(datagram)
	if err != nil {
		return fmt.Errorf("error retrieving in-counter: %v", err)
	}
	if datagram.Counter < prevCounter {
		return fmt.Errorf("replay detected or old datagram: Counter %d is not greater than the last seen in-counter %d", datagram.Counter, prevCounter)
	}

	return nil
}
