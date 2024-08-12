package client_trustlines

import (
    "log"
    "ripple/database/db_trustlines"
    "ripple/handlers"
    "ripple/main"
    "ripple/trustlines" // Import the trustlines package for counter validation
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
        Command:           main.ServerTrustlines_SetTrustline, // Syncing outbound trustline, so use SetTrustline command
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
