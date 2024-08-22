package auth

import (
	"fmt"
	"ripple/types"
	"ripple/database"
)

// validateAndIncrementClientCounter checks if the datagram's counter is valid by comparing it to the last known counter for client connections.
// If valid, it sets the counter to the value in the datagram to prevent replay attacks.
func validateAndIncrementClientCounter(datagram *types.Datagram) error {
	prevCounter, err := database.GetCounter(datagram)
	if err != nil {
		return fmt.Errorf("error retrieving counter: %v", err)
	}
	if datagram.Counter <= prevCounter {
		return fmt.Errorf("replay detected or old datagram: Counter %d is not greater than the last seen counter %d", datagram.Counter, prevCounter)
	}
	if err := database.SetCounter(datagram); err != nil {
		return fmt.Errorf("failed to set counter: %v", err)
	}
	return nil
}

// validateAndIncrementServerCounter checks if the datagram's counter is valid by comparing it to the last known counter for server connections.
// If valid, it sets the counter to the value in the datagram to prevent replay attacks.
func validateAndIncrementServerCounter(datagram *types.Datagram) error {
	prevCounter, err := database.GetCounterIn(datagram)
	if err != nil {
		return fmt.Errorf("error retrieving in-counter: %v", err)
	}
	if datagram.Counter <= prevCounter {
		return fmt.Errorf("replay detected or old datagram: Counter %d is not greater than the last seen in-counter %d", datagram.Counter, prevCounter)
	}
	if err := database.SetCounterIn(datagram); err != nil {
		return fmt.Errorf("failed to set in-counter: %v", err)
	}
	return nil
}

// GetAndIncrementCounterOut retrieves the current counter_out, increments it, and updates the database.
// It returns the counter value before it was incremented.
func GetAndIncrementCounterOut(username, peerServerAddress, peerUsername string) (uint32, error) {
    // Retrieve the current value of counter_out from the database.
    counterOut, err := database.GetCounterOut(username, peerServerAddress, peerUsername)
    if err != nil {
        return 0, err  // Return error if unable to fetch the counter.
    }

    // Increment the counter and update it in the database within the same function call.
    if err := database.SetCounterOut(username, peerServerAddress, peerUsername, counterOut + 1); err != nil {
        return 0, err  // Return error if unable to update the counter.
    }

    // Return the original counter value that was fetched.
    return counterOut, nil
}
