package trustlines

import (
    "fmt"
    "log"
    "ripple/main"
)

// ValidateCounter checks if the provided counter is greater than the stored counter value for counter.
func ValidateCounter(datagram *main.Datagram) error {
    // Retrieve the stored counter value
    prevCounter, err := GetCounter(datagram)
    if err != nil {
        return fmt.Errorf("error getting stored counter for user %s: %v", datagram.Username, err)
    }

    // Check if the incoming counter is valid (greater than the stored counter)
    if datagram.Counter <= prevCounter {
        log.Printf("Received counter (%d) is not greater than stored counter (%d) for user %s. Potential replay attack.", datagram.Counter, prevCounter, datagram.Username)
        return fmt.Errorf("counter validation failed for user %s", datagram.Username)
    }

    return nil
}

// ValidateCounterIn checks if the provided counter is greater than the stored counter_in value.
func ValidateCounterIn(datagram *main.Datagram) error {
    // Retrieve the stored counter_in value
    prevCounterIn, err := GetCounterIn(datagram)
    if err != nil {
        return fmt.Errorf("error getting stored counter_in for user %s: %v", datagram.Username, err)
    }

    // Check if the incoming counter is valid (greater than the stored counter_in)
    if datagram.Counter <= prevCounterIn {
        log.Printf("Received counter (%d) is not greater than stored counter_in (%d) for user %s. Potential replay attack.", datagram.Counter, prevCounterIn, datagram.Username)
        return fmt.Errorf("counter_in validation failed for user %s", datagram.Username)
    }

    return nil
}
