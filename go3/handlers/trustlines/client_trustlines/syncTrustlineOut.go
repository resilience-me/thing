package client_trustlines

import (
    "log"
    "ripple/database/db_trustlines"
    "ripple/main"
    "ripple/trustlines"
)

// SyncTrustlineOut handles the client request to sync the outbound trustline to the peer server.
func SyncTrustlineOut(session main.Session) {
    datagram := session.Datagram

    // Validate the counter using the ValidateCounter function from trustlines package
    if err := trustlines.ValidateCounter(datagram); err != nil {
        log.Printf("Counter validation failed for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Received counter is not valid.", session.Conn)
        return
    }

    // Retrieve the syncCounter and sync status
    syncCounter, isSynced, err := trustlines.GetSyncStatus(datagram)
    if err != nil {
        log.Printf("Failed to retrieve sync status for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to retrieve sync status.", session.Conn)
        return
    }

    // Initialize the datagram
    dg, err := trustlines.InitializeDatagram(datagram)
    if err != nil {
        log.Printf("Error initializing datagram in SyncTrustlineOut for user %s: %v", datagram.Username, err)
        return
    }

    if isSynced {
        // Trustline is already synced, so prepare a SetTimestamp command
        dg.Command = main.ServerTrustlines_SetTimestamp
    } else {
        // Trustline is not synced, proceed with sending the trustline
        trustline, err := db_trustlines.GetTrustlineOut(datagram)
        if err != nil {
            log.Printf("Error getting trustline for user %s: %v", datagram.Username, err)
            main.SendErrorResponse("Failed to retrieve trustline.", session.Conn)
            return
        }
        dg.Command = main.ServerTrustlines_SetTrustline
        binary.BigEndian.PutUint32(dg.Arguments[:4], trustline)
        binary.BigEndian.PutUint32(dg.Arguments[4:8], syncCounter)
    }

    // Send the prepared datagram
    if err := handlers.SignAndSendDatagram(session, dg); err != nil {
        log.Printf("Failed to send datagram for user %s: %v", datagram.Username, err)
        return
    }

    // Update the client-side counter value after processing the datagram
    if err := db_trustlines.SetCounter(datagram, datagram.Counter); err != nil {
        log.Printf("Error updating counter for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to update counter.", session.Conn)
        return
    }

    // Send success response to the client
    if err := main.SendSuccessResponse([]byte("Outbound trustline sync request processed successfully."), session.Conn); err != nil {
        log.Printf("Failed to send success response to user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("SyncTrustline command processed successfully for user %s to peer %s.", datagram.Username, datagram.PeerUsername)
}
