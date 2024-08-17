package main

import (
	"fmt"
)

// validateClientCounter checks if the datagram's counter is valid by comparing it to the last known counter for client connections.
// If valid, it sets the counter to the value in the datagram to prevent replay attacks.
func validateClientCounter(datagram *Datagram) error {
	prevCounter, err := GetCounter(datagram)
	if err != nil {
		return fmt.Errorf("error retrieving counter: %v", err)
	}
	if datagram.Counter < prevCounter {
		return fmt.Errorf("replay detected or old datagram: Counter %d is not greater than the last seen counter %d", datagram.Counter, prevCounter)
	}

	return nil
}

// validateServerCounter checks if the datagram's counter is valid by comparing it to the last known counter for server connections.
// If valid, it sets the counter to the value in the datagram to prevent replay attacks.
func validateServerCounter(datagram *Datagram) error {
	prevCounter, err := GetCounterIn(datagram)
	if err != nil {
		return fmt.Errorf("error retrieving in-counter: %v", err)
	}
	if datagram.Counter < prevCounter {
		return fmt.Errorf("replay detected or old datagram: Counter %d is not greater than the last seen in-counter %d", datagram.Counter, prevCounter)
	}

	return nil
}

// Validate the counter based on its type (client or server)
func ValidateCounter(datagram *Datagram) error {
	if datagram.Command&0x80 == 0 { // Server session if MSB is 0
		return validateServerCounter(dg)
	}
	return validateClientCounter(dg) // Client session if MSB is 1
}

func UpdateClientCounter(datagram *Datagram) error {
	prevCounter, err := GetCounter(datagram)
	if err != nil {
		return fmt.Errorf("error retrieving counter: %v", err)
	}
	if datagram.Counter == prevCounter {
		return nil
	}
	if err := SetCounter(datagram); err != nil {
		return fmt.Errorf("failed to set counter: %v", err)
	}
	return nil
}

func UpdateServerCounter(datagram *Datagram) error {
	prevCounter, err := GetCounterIn(datagram)
	if err != nil {
		return fmt.Errorf("error retrieving in-counter: %v", err)
	}
	if datagram.Counter == prevCounter {
		return nil
	}
	if err := SetCounterIn(datagram); err != nil {
		return fmt.Errorf("failed to set in-counter: %v", err)
	}
	return nil
}
