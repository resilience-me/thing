package db_trustlines

import (
    "ripple/main"
)

// GetAndIncrementCounterOut retrieves the current counter_out, increments it, and updates the database.
// It returns the counter value before it was incremented.
func GetAndIncrementCounterOut(datagram *main.Datagram) (uint32, error) {
    // Retrieve the current value of counter_out from the database.
    counterOut, err := GetCounterOut(datagram)
    if err != nil {
        return 0, err  // Return error if unable to fetch the counter.
    }

    // Increment the counter and update it in the database within the same function call.
    if err := SetCounterOut(datagram, counterOut + 1); err != nil {
        return 0, err  // Return error if unable to update the counter.
    }

    // Return the original counter value that was fetched.
    return counterOut, nil
}

// IncrementSyncCounter retrieves the current sync_counter, increments it, and updates the database.
// It returns an error if something goes wrong during the process.
func IncrementSyncCounter(datagram *main.Datagram) error {
    // Retrieve the current value of sync_counter from the database.
    syncCounter, err := GetSyncCounter(datagram)
    if err != nil {
        return err  // Return error if unable to fetch the sync_counter.
    }

    // Increment the counter and update it in the database within the same function call.
    if err := SetSyncCounter(datagram, syncCounter + 1); err != nil {
        return err  // Return error if unable to update the sync_counter.
    }

    // No need to return any value; just indicate success by returning nil.
    return nil
}
