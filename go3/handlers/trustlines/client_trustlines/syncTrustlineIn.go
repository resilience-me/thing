package client_trustlines

import (
    "log"
    "ripple/main"
    "ripple/handlers"
    "ripple/trustlines"             // Import the trustlines package for counter validation
    "ripple/database/db_trustlines" // Handles database-related operations
)

// SyncTrustlineIn handles the client request to sync the inbound trustline from the peer server.
func SyncTrustlineIn(session main.Session) {
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

    // Create the datagram to request the trustline from the peer
    dg := main.Datagram{
        Command:           main.ServerTrustlines_GetTrustline,
        Username:          datagram.PeerUsername,      // Switch places: this is the peer's username
        PeerUsername:      datagram.Username,          // Switch places: this is your server's username
        PeerServerAddress: serverAddress,              // Your server's address
        Counter:           counterOut,                 // Use the incremented counter_out value
    }

    // Send the GetTrustline command to the peer server
    if err := handlers.SignAndSendDatagram(session, &dg); err != nil {
        log.Printf("Failed to send GetTrustline command for user %s to peer %s: %v", datagram.Username, datagram.PeerUsername, err)
        main.SendErrorResponse("Failed to send GetTrustline command.", session.Conn)
        return
    }

    // Update the client-side counter value after sending the datagram
    if err := db_trustlines.SetCounter(datagram, datagram.Counter); err != nil {
        log.Printf("Error updating counter for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to update counter.", session.Conn)
        return
    }

    // Send success response to the client
    if err := main.SendSuccessResponse([]byte("Trustline sync request sent successfully."), session.Conn); err != nil {
        log.Printf("Failed to send success response to user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("GetTrustline command sent successfully for user %s to peer %s.", datagram.Username, datagram.PeerUsername)
}
