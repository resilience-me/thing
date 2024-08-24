package server_trustlines

import (
    "log"
    "encoding/binary"

    "ripple/comm"
    "ripple/types"
    "ripple/handlers"
    "ripple/handlers/trustlines"
    "ripple/database/db_trustlines"
    "ripple/commands"
)

// GetTrustline handles the request to get the current trustline amount from another server
func GetTrustline(session types.Session) {
    datagram := session.Datagram

    // Retrieve the syncCounter and local sync status
    syncCounter, isSyncedLocally, err := trustlines.GetSyncStatus(datagram)
    if err != nil {
        log.Printf("Failed to retrieve sync status in GetTrustline for user %s: %v", datagram.Username, err)
        return
    }

    // Prepare the datagram
    dg, err := handlers.PrepareDatagramResponse(datagram)
    if err != nil {
        log.Printf("Error preparing datagram in GetTrustline for user %s: %v", datagram.Username, err)
        return
    }

    // Extract sync_in value from the datagram's Arguments[0:4]
    syncIn := types.BytesToUint32(datagram.Arguments[:4])

    if syncIn < syncCounter {
        // The peer is not synced, prepare to send trustline data to synchronize
        dg.Command = commands.ServerTrustlines_SetTrustline

        trustline, err := db_trustlines.GetTrustlineOutFromDatagram(session.Datagram)
        if err != nil {
            log.Printf("Error getting trustline for user %s in GetTrustline: %v", session.Datagram.Username, err)
            return
        }
    
        binary.BigEndian.PutUint32(dg.Arguments[:4], trustline)
        binary.BigEndian.PutUint32(dg.Arguments[4:8], syncCounter)
    } else {
        // Use the SetTimestamp command to the peer to acknowledge synchronization
        dg.Command = commands.ServerTrustlines_SetTimestamp
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
    if err := comm.SignAndSendDatagram(dg, datagram.PeerServerAddress); err != nil {
        log.Printf("Failed to sign and send datagram in GetTrustline for user %s: %v", session.Datagram.Username, err)
        return
    }

    log.Printf("GetTrustline operation completed successfully for user %s.", datagram.Username)
}
