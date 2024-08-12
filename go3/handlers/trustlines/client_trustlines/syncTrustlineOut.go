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

    // Retrieve the sync status
    isSynced, err := trustlines.GetSyncStatus(datagram)
    if err != nil {
        log.Printf("Failed to retrieve sync status for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to retrieve sync status.", session.Conn)
        return
    }

    // Check if the trustline is already synced
    if isSynced {
        // Trustline is already synced, so send a SetTimestamp command instead
        sendSyncTimestamp(session)
    } else {
        // Trustline is not synced, proceed with sending the trustline
        sendTrustline(session, syncCounter)
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

func sendSyncTimestamp(session main.Session) {
    datagram := session.Datagram

    // Retrieve and increment the counter_out value
    counterOut, err := db_trustlines.GetAndIncrementCounterOut(datagram)
    if err != nil {
        log.Printf("Error handling counter_out for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to update counter_out.", session.Conn)
        return
    }

    // Create and send the SetTimestamp command datagram
    dg := main.Datagram{
        Command:           main.ServerTrustlines_SetTimestamp,
        Username:          datagram.Username,
        PeerUsername:      datagram.PeerUsername,
        PeerServerAddress: datagram.PeerServerAddress,
        Counter:           counterOut,
    }

    if err := handlers.SignAndSendDatagram(session, &dg); err != nil {
        log.Printf("Failed to send SetTimestamp command for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to send SetTimestamp command.", session.Conn)
        return
    }

    log.Printf("SetTimestamp command sent successfully for user %s to peer %s.", datagram.Username, datagram.PeerUsername)
}

func sendTrustline(session main.Session, syncCounter uint32) {
    datagram := session.Datagram

    // Retrieve and increment the counter_out value
    counterOut, err := db_trustlines.GetAndIncrementCounterOut(datagram)
    if err != nil {
        log.Printf("Error handling counter_out for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to update counter_out.", session.Conn)
        return
    }

    // Retrieve the trustline amount to be sent
    trustline, err := db_trustlines.GetTrustlineOut(datagram)
    if err != nil {
        log.Printf("Error getting trustline for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to retrieve trustline.", session.Conn)
        return
    }

    // Create and send the SetTrustline command datagram
    dg := main.Datagram{
        Command:           main.ServerTrustlines_SetTrustline,
        Username:          datagram.Username,
        PeerUsername:      datagram.PeerUsername,
        PeerServerAddress: datagram.PeerServerAddress,
        Counter:           counterOut,
    }
    binary.BigEndian.PutUint32(dg.Arguments[:4], trustline)
    binary.BigEndian.PutUint32(dg.Arguments[4:8], syncCounter)

    if err := handlers.SignAndSendDatagram(session, &dg); err != nil {
        log.Printf("Failed to send SetTrustline command for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to send SetTrustline command.", session.Conn)
        return
    }

    log.Printf("SetTrustline command sent successfully for user %s to peer %s.", datagram.Username, datagram.PeerUsername)
}
