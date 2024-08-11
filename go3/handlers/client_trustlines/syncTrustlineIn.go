package client_trustlines

import (
    "log"
    "ripple/main"
    "ripple/handlers"
    "ripple/database/db_trustlines"
)

// SyncTrustlineIn handles the client request to sync the inbound trustline from the peer server.
func SyncTrustlineIn(session main.Session) {
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

    // Create the datagram to request the trustline from the peer
    dg := main.Datagram{
        Command:           main.ServerTrustlines_GetTrustline,
        Username:          datagram.Username,         // Your server's username
        PeerUsername:      datagram.PeerUsername,     // The peer's username
        PeerServerAddress: datagram.PeerServerAddress, // The peer's server address
        Counter:           datagram.Counter,          // Use the session's counter
    }

    // Send the GetTrustline command to the peer server
    if err := handlers.SignAndSendDatagram(session, &dg); err != nil {
        log.Printf("Failed to send GetTrustline command for user %s to peer %s: %v", datagram.Username, datagram.PeerUsername, err)
        return
    }

    // Update the counter value after sending the datagram
    if err := db_trustlines.SetCounter(datagram, datagram.Counter); err != nil {
        log.Printf("Error updating counter for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to update counter.", session.Conn)
        return
    }

    log.Printf("GetTrustline command sent successfully for user %s to peer %s.", datagram.Username, datagram.PeerUsername)
}
