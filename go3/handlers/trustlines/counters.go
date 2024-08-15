package trustlines

import (
    "fmt"
    "ripple/main"
    "ripple/database/db_trustlines"
)

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

// GetSyncStatus retrieves the syncCounter and syncOut values and returns the syncCounter and whether they are equal.
func GetSyncStatus(datagram *main.Datagram) (uint32, bool, error) {
    // Retrieve the current syncCounter value
    syncCounter, err := db_trustlines.GetSyncCounter(datagram)
    if err != nil {
        return 0, false, fmt.Errorf("Error getting syncCounter for user %s: %v", datagram.Username, err)
    }

    // Retrieve the current syncOut value
    syncOut, err := db_trustlines.GetSyncOut(datagram)
    if err != nil {
        return 0, false, fmt.Errorf("Error getting syncOut for user %s: %v", datagram.Username, err)
    }

    return syncCounter, syncOut == syncCounter, nil
}
