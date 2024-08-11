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
    // Retrieve the current sync_counter value
    syncCounter, err := db_trustlines.GetSyncCounter(session.Datagram)
    if err != nil {
        log.Printf("Error getting sync_counter for user %s: %v", session.Datagram.Username, err)
        return
    }

    // Retrieve the current sync_out value
    syncOut, err := db_trustlines.GetSyncOut(session.Datagram)
    if err != nil {
        log.Printf("Error getting sync_out for user %s: %v", session.Datagram.Username, err)
        return
    }

    // Retrieve and increment the counter_out value
    counterOut, err := db_trustlines.GetAndIncrementCounterOut(session.Datagram)
    if err != nil {
        log.Printf("Error handling counter_out for user %s: %v", session.Datagram.Username, err)
        return
    }

    // Initialize common Datagram fields
    dg := main.Datagram{
        Username:          session.Datagram.PeerUsername,
        PeerUsername:      session.Datagram.Username,
        PeerServerAddress: main.GetServerAddress(),
        Counter:           counterOut,
    }

    if syncCounter == syncOut {
        sendSyncTimestamp(session, &dg)
    } else {
        sendTrustline(session, &dg, syncCounter)
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
