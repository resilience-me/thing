package client_trustlines

import (
    "encoding/binary"
    "log"

    "ripple/commands"
    "ripple/comm"
    "ripple/database/db_trustlines"
    "ripple/handlers"
    "ripple/types"
    "ripple/handlers/trustlines"
)

// SyncTrustlineOut handles the client request to sync the outbound trustline to the peer server.
func SyncTrustlineOut(session types.Session) {
    datagram := session.Datagram

    // Prepare the datagram
    dgOut, err := handlers.PrepareDatagramResponse(datagram)
    if err != nil {
        log.Printf("Error preparing datagram for user %s: %v", datagram.Username, err)
        comm.SendErrorResponse(session.Addr, "Error preparing datagram.")
        return
    }

    // Retrieve the syncCounter and sync status
    syncCounter, isSynced, err := trustlines.GetSyncStatus(datagram)
    if err != nil {
        log.Printf("Failed to retrieve sync status in SyncTrustlineOut for user %s: %v", datagram.Username, err)
        comm.SendErrorResponse(session.Addr, "Failed to retrieve sync status.")
        return
    }

    if isSynced {
        // Trustline is already synced, so prepare a SetTimestamp command
        dgOut.Command = commands.ServerTrustlines_SetTimestamp
    } else {
        // Trustline is not synced, proceed with sending the trustline
        trustline, err := db_trustlines.GetTrustlineOut(datagram)
        if err != nil {
            log.Printf("Error getting trustline for user %s in SyncTrustlineOut: %v", datagram.Username, err)
            comm.SendErrorResponse(session.Addr, "Failed to retrieve trustline.")
            return
        }
        dgOut.Command = commands.ServerTrustlines_SetTrustline
        binary.BigEndian.PutUint32(dgOut.Arguments[:4], trustline)
        binary.BigEndian.PutUint32(dgOut.Arguments[4:8], syncCounter)
    }

    // Send the prepared datagram
    if err := comm.SignAndSendDatagram(dgOut, datagram.PeerServerAddress); err != nil {
        log.Printf("Failed to send datagram in SyncTrustlineOut for user %s: %v", datagram.Username, err)
        return
    }

    // Send success response to the client
    if err := comm.SendSuccessResponse(session.Addr, []byte("Outbound trustline sync request processed successfully.")); err != nil {
        log.Printf("Failed to send success response in SyncTrustlineOut for user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("SyncTrustline command processed successfully for user %s to peer %s.", datagram.Username, datagram.PeerUsername)
}
