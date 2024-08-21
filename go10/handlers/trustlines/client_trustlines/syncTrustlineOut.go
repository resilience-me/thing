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

    // Retrieve the syncCounter and sync status
    syncCounter, isSynced, err := trustlines.GetSyncStatus(datagram)
    if err != nil {
        log.Printf("Failed to retrieve sync status in SyncTrustlineOut for user %s: %v", datagram.Username, err)
        comm.SendErrorResponse("Failed to retrieve sync status.", session.Conn)
        return
    }

    // Retrieve and increment the counter_out value
    counterOut, err := db_trustlines.GetAndIncrementCounterOut(datagram)
    if err != nil {
        log.Printf("Error handling counter_out for user %s: %v", datagram.Username, err)
        comm.SendErrorResponse("Failed to update counter_out.", session.Conn)
        return
    }

    dgOut := types.NewDatagram(datagram.PeerUsername, datagram.Username, counterOut)

    if isSynced {
        // Trustline is already synced, so prepare a SetTimestamp command
        dgOut.Command = main.ServerTrustlines_SetTimestamp
    } else {
        // Trustline is not synced, proceed with sending the trustline
        trustline, err := db_trustlines.GetTrustlineOut(datagram)
        if err != nil {
            log.Printf("Error getting trustline for user %s in SyncTrustlineOut: %v", datagram.Username, err)
            comm.SendErrorResponse("Failed to retrieve trustline.", session.Conn)
            return
        }
        dgOut.Command = main.ServerTrustlines_SetTrustline
        binary.BigEndian.PutUint32(dgOut.Arguments[:4], trustline)
        binary.BigEndian.PutUint32(dgOut.Arguments[4:8], syncCounter)
    }

    // Send the prepared datagram
    if err := comm.SignAndSendDatagram(session, dgOut); err != nil {
        log.Printf("Failed to send datagram in SyncTrustlineOut for user %s: %v", datagram.Username, err)
        return
    }

    // Send success response to the client
    if err := comm.SendSuccessResponse([]byte("Outbound trustline sync request processed successfully."), session.Addr); err != nil {
        log.Printf("Failed to send success response in SyncTrustlineOut for user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("SyncTrustline command processed successfully for user %s to peer %s.", datagram.Username, datagram.PeerUsername)
}