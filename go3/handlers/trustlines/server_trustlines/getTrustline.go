package server_trustlines

import (
    "log"
    "ripple/main"
    "ripple/handlers"
    "ripple/handlers/trustlines"
    "ripple/database/db_trustlines"
    "ripple/database/db_server"
)

// GetTrustline handles the request to get the current trustline amount from another server
func GetTrustline(session main.Session) {
    datagram := session.Datagram

    // Validate the counter_in to ensure the request is not a replay
    if err := db_server.ValidateCounterIn(datagram); err != nil {
        log.Printf("Counter_in validation failed for user %s: %v", datagram.Username, err)
        return
    }

    // Retrieve the syncCounter and local sync status
    syncCounter, isSyncedLocally, err := trustlines.GetSyncStatus(datagram)
    if err != nil {
        log.Printf("Failed to retrieve sync status in GetTrustline for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to retrieve sync status.", session.Conn)
        return
    }

    // Extract sync_in value from the datagram's Arguments[0:4]
    syncIn := main.BytesToUint32(datagram.Arguments[:4])

    // Initialize the datagram
    dg, err := trustlines.InitializeDatagram(datagram)
    if err != nil {
        log.Printf("Error initializing datagram in GetTrustline for user %s: %v", datagram.Username, err)
        return
    }

    if syncIn < syncCounter {
        // The peer is not synced, prepare to send trustline data to synchronize
        dg.Command = main.ServerTrustlines_SetTrustline

        trustline, err := db_trustlines.GetTrustlineOut(session.Datagram)
        if err != nil {
            log.Printf("Error getting trustline for user %s in GetTrustline: %v", session.Datagram.Username, err)
            return
        }
    
        binary.BigEndian.PutUint32(dg.Arguments[:4], trustline)
        binary.BigEndian.PutUint32(dg.Arguments[4:8], syncCounter)
    } else {
        // Use the SetTimestamp command to the peer to acknowledge synchronization
        dg.Command = main.ServerTrustlines_SetTimestamp
        if !isSyncedLocally {
            // The peer is synced, but the local server is not aware
            // Update the local sync_out to match the sync_counter
            if err := db_trustlines.SetSyncOut(datagram, syncCounter); err != nil {
                log.Printf("Error updating sync_out in GetTrustline for user %s: %v", datagram.Username, err)
                return
            }
        }
    }

    // Send the prepared datagram
    if err := handlers.SignAndSendDatagram(session, dg); err != nil {
        log.Printf("Failed to sign and send datagram in GetTrustline for user %s: %v", session.Datagram.Username, err)
        return
    }
    log.Printf("Datagram command %d sent successfully in GetTrustline for user %s.", dg.Command, session.Datagram.Username)

    // Update the counter_in after successfully processing the request
    if err := db_server.SetCounterIn(datagram, datagram.Counter); err != nil {
        log.Printf("Error updating counter_in in GetTrustline for user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("GetTrustline operation completed successfully for user %s.", datagram.Username)
}
