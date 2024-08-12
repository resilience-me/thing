package server_trustlines

import (
    "encoding/binary"
    "log"
    "ripple/main"
    "ripple/handlers"
    "ripple/handlers/trustlines"
    "ripple/database/db_trustlines"
)

// GetTrustline handles the request to get the current trustline amount from another server
func GetTrustline(session main.Session) {
    datagram := session.Datagram

    // Validate the counter_in to ensure the request is not a replay
    if err := trustlines.ValidateCounterIn(datagram); err != nil {
        log.Printf("Counter_in validation failed for user %s: %v", datagram.Username, err)
        return
    }

    // Retrieve the current sync_counter value
    syncCounter, err := db_trustlines.GetSyncCounter(datagram)
    if err != nil {
        log.Printf("Error getting sync_counter for user %s: %v", datagram.Username, err)
        return
    }

    // Extract sync_in value from the datagram's Arguments[0:4]
    syncIn := main.BytesToUint32(datagram.Arguments[:4])

    // Retrieve the current sync_out value
    syncOut, err := db_trustlines.GetSyncOut(datagram)
    if err != nil {
        log.Printf("Error getting sync_out for user %s: %v", datagram.Username, err)
        return
    }

    // Retrieve and increment the counter_out value
    counterOut, err := trustlines.GetAndIncrementCounterOut(datagram)
    if err != nil {
        log.Printf("Error handling counter_out for user %s: %v", datagram.Username, err)
        return
    }

    // Initialize common Datagram fields for response
    dg := main.Datagram{
        Username:          datagram.PeerUsername,
        PeerUsername:      datagram.Username,
        PeerServerAddress: main.GetServerAddress(),
        Counter:           counterOut,
    }

    // Logic to determine the correct response
    if syncIn < syncCounter {
        // The peer is not synced, prepare to send trustline data to synchronize
        dg.Command = main.ServerTrustlines_SetTrustline

        trustline, err := db_trustlines.GetTrustlineOut(session.Datagram)
        if err != nil {
            log.Printf("Error getting trustline for user %s: %v", session.Datagram.Username, err)
            return
        }
    
        binary.BigEndian.PutUint32(dg.Arguments[:4], trustline)
        binary.BigEndian.PutUint32(dg.Arguments[4:8], syncCounter)
    } else {
        // Use the SetTimestamp command to the peer to acknowledge synchronization
        dg.Command = main.ServerTrustlines_SetTimestamp
        if syncOut < syncCounter {
            // The peer is synced, but the local server is not aware
            // Update the local sync_out to match the sync_counter
            if err := db_trustlines.SetSyncOut(datagram, syncCounter); err != nil {
                log.Printf("Error updating sync_out for user %s: %v", datagram.Username, err)
                return
            }
        }
    }

    if err := handlers.SignAndSendDatagram(session, dg); err != nil {
        log.Printf("Failed to sign and send datagram for user %s: %v", session.Datagram.Username, err)
        return
    }
    log.Printf("SetTrustlineSyncTimestamp command sent successfully for user %s.", session.Datagram.Username)

    // Update the counter_in after successfully processing the request
    if err := db_trustlines.SetCounterIn(datagram, datagram.Counter); err != nil {
        log.Printf("Error updating counter_in for user %s: %v", datagram.Username, err)
        return
    }

}
