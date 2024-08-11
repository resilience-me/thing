package client_trustlines

import (
    "log"
    "ripple/main"
    "ripple/handlers"
)

// SyncTrustlineIn handles the client request to sync the inbound trustline from the peer server.
func SyncTrustlineIn(session main.Session) {
    // Create the datagram to request the trustline from the peer
    dg := main.Datagram{
        Command:           main.ServerTrustlines_GetTrustline,
        Username:          session.Datagram.Username,  // Your server's username
        PeerUsername:      session.Datagram.PeerUsername,  // The peer's username
        PeerServerAddress: session.Datagram.PeerServerAddress,  // The peer's server address
        Counter:           session.Datagram.Counter,  // Use the session's counter
    }

    // Send the GetTrustline command to the peer server
    if err := handlers.SignAndSendDatagram(session, &dg); err != nil {
        log.Printf("Failed to send GetTrustline command for user %s to peer %s: %v", session.Datagram.Username, session.Datagram.PeerUsername, err)
        return
    }

    log.Printf("GetTrustline command sent successfully for user %s to peer %s.", session.Datagram.Username, session.Datagram.PeerUsername)
}
