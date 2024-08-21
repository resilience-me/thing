package client_trustlines

import (
    "log"
    "encoding/binary"
    "ripple/main"
    "ripple/config"
    "ripple/comm"
    "ripple/types"
    "ripple/database/db_trustlines"
)

// SyncTrustlineIn handles the client request to sync the inbound trustline from the peer server.
func SyncTrustlineIn(session main.Session) {
    datagram := session.Datagram

    // Retrieve and increment the counter_out value
    counterOut, err := db_trustlines.GetAndIncrementCounterOut(datagram)
    if err != nil {
        log.Printf("Error handling counter_out for user %s: %v", datagram.Username, err)
        comm.SendErrorResponse("Failed to update counter_out.", session.Conn)
        return
    }

    // Retrieve the current sync_in value
    syncIn, err := db_trustlines.GetSyncIn(datagram)
    if err != nil {
        log.Printf("Error getting sync_in for user %s: %v", datagram.Username, err)
        comm.SendErrorResponse("Failed to read sync_in value.", session.Conn)
        return
    }

    // Fetch the server's address
    serverAddress := config.GetServerAddress()

    // Create the datagram to request the trustline from the peer, including the sync_in value
    dg := types.Datagram{
        Command:           main.ServerTrustlines_GetTrustline,
        Username:          datagram.PeerUsername,      // Switch places: this is the peer's username
        PeerUsername:      datagram.Username,          // Switch places: this is your server's username
        PeerServerAddress: serverAddress,              // Your server's address
        Counter:           counterOut,                 // Use the incremented counter_out value
    }

    // Include the sync_in value in the datagram's Arguments[0:4]
    binary.BigEndian.PutUint32(dg.Arguments[0:4], syncIn)

    // Send the GetTrustline command to the peer server
    if err := comm.SignAndSendDatagram(session, &dg); err != nil {
        log.Printf("Failed to send GetTrustline command for user %s to peer %s: %v", datagram.Username, datagram.PeerUsername, err)
        comm.SendErrorResponse("Failed to send GetTrustline command.", session.Conn)
        return
    }

    // Send success response to the client
    if err := comm.SendSuccessResponse([]byte("Trustline sync request sent successfully."), session.Conn); err != nil {
        log.Printf("Failed to send success response to user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("GetTrustline command sent successfully for user %s to peer %s.", datagram.Username, datagram.PeerUsername)
}
