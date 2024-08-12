package trustlines

import (
    "fmt"
    "ripple/main"
    "ripple/database/db_trustlines"
)

// ValidateCounter checks if the provided counter is greater than the stored counter value for counter.
func ValidateCounter(datagram *main.Datagram) error {
    // Retrieve the stored counter value
    prevCounter, err := db_trustlines.GetCounter(datagram)
    if err != nil {
        return fmt.Errorf("error getting stored counter for user %s: %v", datagram.Username, err)
    }

    // Check if the incoming counter is valid (greater than the stored counter)
    if datagram.Counter <= prevCounter {
        return fmt.Errorf("counter validation failed for user %s", datagram.Username)
    }

    return nil
}

// ValidateCounterIn checks if the provided counter is greater than the stored counter_in value.
func ValidateCounterIn(datagram *main.Datagram) error {
    // Retrieve the stored counter_in value
    prevCounterIn, err := db_trustlines.GetCounterIn(datagram)
    if err != nil {
        return fmt.Errorf("error getting stored counter_in for user %s: %v", datagram.Username, err)
    }

    // Check if the incoming counter is valid (greater than the stored counter_in)
    if datagram.Counter <= prevCounterIn {
        return fmt.Errorf("counter_in validation failed for user %s", datagram.Username)
    }

    return nil
}

// GetAndIncrementCounterOut retrieves the current counter_out, increments it, and updates the database.
// It returns the counter value before it was incremented.
func GetAndIncrementCounterOut(datagram *main.Datagram) (uint32, error) {
    // Retrieve the current value of counter_out from the database.
    counterOut, err := db_trustlines.GetCounterOut(datagram)
    if err != nil {
        return 0, err  // Return error if unable to fetch the counter.
    }

    // Increment the counter and update it in the database within the same function call.
    if err := db_trustlines.SetCounterOut(datagram, counterOut + 1); err != nil {
        return 0, err  // Return error if unable to update the counter.
    }

    // Return the original counter value that was fetched.
    return counterOut, nil
}

// IncrementSyncCounter retrieves the current sync_counter, increments it, and updates the database.
// It returns an error if something goes wrong during the process.
func IncrementSyncCounter(datagram *main.Datagram) error {
    // Retrieve the current value of sync_counter from the database.
    syncCounter, err := db_trustlines.GetSyncCounter(datagram)
    if err != nil {
        return err  // Return error if unable to fetch the sync_counter.
    }

    // Increment the counter and update it in the database within the same function call.
    if err := db_trustlines.SetSyncCounter(datagram, syncCounter + 1); err != nil {
        return err  // Return error if unable to update the sync_counter.
    }

    // No need to return any value; just indicate success by returning nil.
    return nil
}

// GetSyncStatus retrieves the syncCounter and syncOut values and returns whether they are equal.
func GetSyncStatus(datagram *main.Datagram) (bool, error) {
    // Retrieve the current syncCounter value
    syncCounter, err := db_trustlines.GetSyncCounter(datagram)
    if err != nil {
        return false, fmt.Errorf("Error getting syncCounter for user %s: %v", datagram.Username, err)
    }

    // Retrieve the current syncOut value
    syncOut, err := db_trustlines.GetSyncOut(datagram)
    if err != nil {
        return false, fmt.Errorf("Error getting syncOut for user %s: %v", datagram.Username, err)
    }

    return syncOut == syncCounter, nil
}
