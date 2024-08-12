package server_trustlines

import (
    "encoding/binary"
    "log"
    "ripple/main"
    "ripple/handlers"
    "ripple/database/db_trustlines"
)

// GetTrustline handles the request to get the current trustline amount from another server
func GetTrustline(session main.Session) {
    datagram := session.Datagram

    // Retrieve and validate the current counter_in value
    prevCounterIn, err := db_trustlines.GetCounterIn(datagram)
    if err != nil {
        log.Printf("Error getting counter_in for user %s: %v", datagram.Username, err)
        return
    }

    // Check if the incoming counter is valid (prevents replay attacks)
    if datagram.Counter <= prevCounterIn {
        log.Printf("Received counter (%d) is not greater than previous counter_in (%d) for user %s. Potential replay attack.", datagram.Counter, prevCounterIn, datagram.Username)
        return
    }

    // Retrieve the current sync_counter value
    syncCounter, err := db_trustlines.GetSyncCounter(datagram)
    if err != nil {
        log.Printf("Error getting sync_counter for user %s: %v", datagram.Username, err)
        return
    }

    // Retrieve the current sync_out value
    syncOut, err := db_trustlines.GetSyncOut(datagram)
    if err != nil {
        log.Printf("Error getting sync_out for user %s: %v", datagram.Username, err)
        return
    }

    // Retrieve and increment the counter_out value
    counterOut, err := db_trustlines.GetAndIncrementCounterOut(datagram)
    if err != nil {
        log.Printf("Error handling counter_out for user %s: %v", datagram.Username, err)
        return
    }

    // Initialize common Datagram fields
    dg := main.Datagram{
        Username:          datagram.PeerUsername,
        PeerUsername:      datagram.Username,
        PeerServerAddress: main.GetServerAddress(),
        Counter:           counterOut,
    }

    // Send the appropriate response based on sync status
    if syncCounter == syncOut {
        sendSyncTimestamp(session, &dg)
    } else {
        sendTrustline(session, &dg, syncCounter)
    }

    // Update the counter_in after successfully processing the request
    if err := db_trustlines.SetCounterIn(datagram, datagram.Counter); err != nil {
        log.Printf("Error updating counter_in for user %s: %v", datagram.Username, err)
        return
    }
}

func sendSyncTimestamp(session main.Session, dg *main.Datagram) {
    dg.Command = main.ServerTrustlines_SetTimestamp

    if err := handlers.SignAndSendDatagram(session, dg); err != nil {
        log.Printf("Failed to sign and send datagram for user %s: %v", session.Datagram.Username, err)
        return
    }
    log.Printf("SetTrustlineSyncTimestamp command sent successfully for user %s.", session.Datagram.Username)
}

func sendTrustline(session main.Session, dg *main.Datagram, syncCounter uint32) {
    trustline, err := db_trustlines.GetTrustlineOut(session.Datagram)
    if err != nil {
        log.Printf("Error getting trustline for user %s: %v", session.Datagram.Username, err)
        return
    }

    dg.Command = main.ServerTrustlines_SetTrustline
    binary.BigEndian.PutUint32(dg.Arguments[:4], trustline)
    binary.BigEndian.PutUint32(dg.Arguments[4:8], syncCounter)

    if err := handlers.SignAndSendDatagram(session, dg); err != nil {
        log.Printf("Failed to sign and send datagram for user %s: %v", session.Datagram.Username, err)
        return
    }
    log.Printf("SetTrustline command sent successfully for user %s.", session.Datagram.Username)
}
