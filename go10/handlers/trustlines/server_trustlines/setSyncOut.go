package server_trustlines

import (
    "log"
    "ripple/types"
    "ripple/database/db_trustlines"
)

// SetSyncOut handles updating the sync_out counter from a received context
func SetSyncOut(session types.Session) {
    datagram := session.Datagram

    // Load the new sync_out value from the Arguments in the Datagram
    syncOut := types.BytesToUint32(datagram.Arguments[:4])

    // Retrieve the previous sync_out value
    prevSyncOut, err := db_trustlines.GetSyncOut(datagram)
    if err != nil {
        log.Printf("Error getting previous sync_out for user %s: %v", datagram.Username, err)
        return
    }

    // Check if the new sync_out is greater than the previous sync_out
    if syncOut <= prevSyncOut {
        log.Printf("Received sync_out (%d) is not greater than previous sync_out (%d) for user %s.",
            syncOut, prevSyncOut, datagram.Username)
        return
    }

    // Write the new sync_out value
    if err := db_trustlines.SetSyncOut(datagram, syncOut); err != nil {
        log.Printf("Error writing sync_out to file for user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("Sync_out updated successfully for user %s.", datagram.Username)
}
