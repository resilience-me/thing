package server_trustlines

import (
    "log"
    "time"
    "ripple/handlers"
    "ripple/commands"
    "ripple/types"
    "ripple/database/db_trustlines"
)

// SetTrustline handles setting or updating a trustline from another server's perspective.
func SetTrustline(session types.Session) {
    datagram := session.Datagram

    // Retrieve the sync_in value using the new getter
    prevSyncIn, err := db_trustlines.GetSyncIn(datagram)
    if err != nil {
        log.Printf("Error getting sync_in for user %s: %v", datagram.Username, err)
        return
    }

    // Retrieve the syncIn counter from the Datagram
    syncInBytes := datagram.Arguments[4:8]
    syncIn := types.BytesToUint32(syncInBytes)

    if syncIn > prevSyncIn {

        // Retrieve the trustline amount from the Datagram
        trustlineAmount := types.BytesToUint32(datagram.Arguments[:4])
    
        // Update the trustline, sync_in, and timestamp
        if err := db_trustlines.SetTrustlineInFromDatagram(datagram, trustlineAmount); err != nil {
            log.Printf("Error writing trustline to file for user %s: %v", datagram.Username, err)
            return
        }
    
        if err := db_trustlines.SetSyncIn(datagram, syncIn); err != nil {
            log.Printf("Error writing sync_in to file for user %s: %v", datagram.Username, err)
            return
        }
    
        log.Printf("Trustline and sync_in updated successfully for user %s.", datagram.Username)
    
        // Prepare and send the datagram to sync the peer's out counter
        if err := handlers.PrepareAndSendDatagram(commands.ServerTrustlines_SetSyncOut, datagram.Username, datagram.PeerServerAddress, datagram.PeerUsername, syncInBytes); err != nil {
            log.Printf("Failed to sign and send datagram for user %s: %v", datagram.Username, err)
            return
        }
    
        log.Printf("Trustline update and datagram sent successfully for user %s.", datagram.Username)        
    } else {
        log.Printf("Sync_in is synchronized with the peer's most recent trustline_out for user %s.", datagram.Username)
    }

    if err := db_trustlines.SetTimestamp(datagram, time.Now().Unix()); err != nil {
        log.Printf("Error writing timestamp to file for user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("Trustline synchronization timestamp updated successfully for user %s.", datagram.Username)
}
