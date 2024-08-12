package server_trustlines

import (
    "log"
    "ripple/main"
    "ripple/trustlines"             // Import the trustlines package for counter validation
    "ripple/database/db_trustlines" // Handles database-related operations
)

// SetSyncOut handles updating the sync_out counter from a received context
func SetSyncOut(session main.Session) {
    datagram := session.Datagram

    // Validate the counter_in using the ValidateCounterIn function from trustlines package
    if err := trustlines.ValidateCounterIn(datagram); err != nil {
        log.Printf("Counter_in validation failed for user %s: %v", datagram.Username, err)
        return
    }

    // Load the new sync_out value from the Arguments in the Datagram
    syncOut := main.BytesToUint32(datagram.Arguments[:4])

    // Retrieve the previous sync_out value
    prevSyncOut, err := db_trustlines.GetSyncOut(datagram)
    if err != nil {
        log.Printf("Error getting previous sync_out for user %s: %v", datagram.Username, err)
        return
    }

    // Check if the new sync_out is greater than the previous sync_out
    if syncOut <= prevSyncOut {
        log.Printf("Received sync_out (%d) is not greater than previous sync_out (%d) for user %s. Potential replay attack.",
            syncOut, prevSyncOut, datagram.Username)
        return
    }

    // Write the new sync_out value
    if err := db_trustlines.SetSyncOut(datagram, syncOut); err != nil {
        log.Printf("Error writing sync_out to file for user %s: %v", datagram.Username, err)
        return
    }

    // After successfully updating sync_out, update the counter_in
    if err := db_trustlines.SetCounterIn(datagram, datagram.Counter); err != nil {
        log.Printf("Error updating counter_in for user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("Sync_out updated successfully for user %s.", datagram.Username)
}
