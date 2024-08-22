package client_trustlines

import (
    "log"
    "encoding/binary"
    "ripple/main"
    "ripple/comm"
    "ripple/handlers"
    "ripple/database/db_trustlines"
    "ripple/commands"
)

// SyncTrustlineIn handles the client request to sync the inbound trustline from the peer server.
func SyncTrustlineIn(session main.Session) {
    datagram := session.Datagram

    // Prepare the datagram
    dgOut, err := handlers.PrepareDatagramResponse(datagram)
    if err != nil {
        log.Printf("Error preparing datagram for user %s: %v", datagram.Username, err)
        comm.SendErrorResponse("Error preparing datagram.", session.Addr)
        return
    }

    // Retrieve the current sync_in value
    syncIn, err := db_trustlines.GetSyncIn(datagram)
    if err != nil {
        log.Printf("Error getting sync_in for user %s: %v", datagram.Username, err)
        comm.SendErrorResponse("Failed to read sync_in value.", session.Addr)
        return
    }

    dgOut.Command = commands.ServerTrustlines_GetTrustline
    // Include the sync_in value in the datagram's Arguments[0:4]
    binary.BigEndian.PutUint32(dgOut.Arguments[0:4], syncIn)

    // Send the GetTrustline command to the peer server
    if err := comm.SignAndSendDatagram(session, dgOut); err != nil {
        log.Printf("Failed to send GetTrustline command for user %s to peer %s: %v", datagram.Username, datagram.PeerUsername, err)
        comm.SendErrorResponse("Failed to send GetTrustline command.", session.Addr)
        return
    }

    // Send success response to the client
    if err := comm.SendSuccessResponse([]byte("Trustline sync request sent successfully."), session.Addr); err != nil {
        log.Printf("Failed to send success response to user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("GetTrustline command sent successfully for user %s to peer %s.", datagram.Username, datagram.PeerUsername)
}
