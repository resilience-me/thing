package client_trustlines

import (
    "log"
    "ripple/database/db_trustlines"
    "ripple/handlers"
    "ripple/main"
)

// SyncTrustlineOut handles the client request to sync the outbound trustline to the peer server.
func SyncTrustlineOut(session main.Session) {
    datagram := session.Datagram

    // Retrieve the previous client-side counter value using the getter
    prevCounter, err := db_trustlines.GetCounter(datagram)
    if err != nil {
        log.Printf("Error getting previous counter for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to read counter file.", session.Conn)
        return
    }

    // Check if the client-side counter is valid (prevents replay attacks)
    if datagram.Counter <= prevCounter {
        log.Printf("Received counter is not greater than previous counter for user %s. Potential replay attack.", datagram.Username)
        main.SendErrorResponse("Received counter is not valid.", session.Conn)
        return
    }

    // Retrieve and increment the counter_out value
    counterOut, err := db_trustlines.GetAndIncrementCounterOut(datagram)
    if err != nil {
        log.Printf("Error handling counter_out for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to update counter_out.", session.Conn)
        return
    }

    // Fetch the server's address
    serverAddress := main.GetServerAddress()

    // Create the datagram to sync the outbound trustline to the peer
    dg := main.Datagram{
        Command:           main.ServerTrustlines_GetTrustline, // Assuming you want to use the same command as SyncTrustlineIn
        Username:          datagram.Username,
        PeerUsername:      datagram.PeerUsername,
        PeerServerAddress: datagram.PeerServerAddress,
        Counter:           counterOut,
    }

    // Send the SyncTrustline command to the peer server
    if err := handlers.SignAndSendDatagram(session, &dg); err != nil {
        log.Printf("Failed to send SyncTrustline command for user %s to peer %s: %v", datagram.Username, datagram.PeerUsername, err)
        main.SendErrorResponse("Failed to send SyncTrustline command.", session.Conn)
        return
    }

    // Update the client-side counter value after sending the datagram
    if err := db_trustlines.SetCounter(datagram, datagram.Counter); err != nil {
        log.Printf("Error updating counter for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to update counter.", session.Conn)
        return
    }

    // Send success response to the client
    if err := main.SendSuccessResponse([]byte("Outbound trustline sync request sent successfully."), session.Conn); err != nil {
        log.Printf("Failed to send success response to user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("SyncTrustline command sent successfully for user %s to peer %s.", datagram.Username, datagram.PeerUsername)
}
