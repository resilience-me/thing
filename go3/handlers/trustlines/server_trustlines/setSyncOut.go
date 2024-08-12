package server_trustlines

import (
    "log"
    "ripple/main"
    "ripple/database/db_trustlines"
)

// SetSyncOut handles updating the sync_out counter from a received context
func SetSyncOut(session main.Session) {
    // Retrieve the previous sync_out value
    prevSyncOut, err := db_trustlines.GetSyncOut(session.Datagram)
    if err != nil {
        log.Printf("Error getting previous sync_out for user %s: %v", session.Datagram.Username, err)
        return
    }

    // Get the new sync_out value from the datagram
    syncOut := session.Datagram.Counter

    // Check if the new sync_out is greater than the previous sync_out
    if syncOut <= prevSyncOut {
        log.Printf("Received sync_out (%d) is not greater than previous sync_out (%d) for user %s. Potential replay attack.",
            syncOut, prevSyncOut, session.Datagram.Username)
        return
    }

    // Write the new sync_out value
    if err := db_trustlines.SetSyncOut(session.Datagram, syncOut); err != nil {
        log.Printf("Error writing sync_out to file for user %s: %v", session.Datagram.Username, err)
        return
    }

    log.Printf("Sync_out updated successfully for user %s.", session.Datagram.Username)
}
