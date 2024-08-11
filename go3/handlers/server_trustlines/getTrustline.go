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
    // Retrieve the current counter value
    counter, err := db_trustlines.GetCounter(session.Datagram)
    if err != nil {
        log.Printf("Error getting counter: %v", err)
        return
    }

    // Retrieve the current sync_out value
    syncOut, err := db_trustlines.GetSyncOut(session.Datagram)
    if err != nil {
        log.Printf("Error getting sync_out: %v", err)
        return
    }

    if counter == syncOut {
        sendSyncTimestamp(session)
    } else {
        sendTrustline(session, counter)
    }
}

func sendSyncTimestamp(session main.Session) {
    syncCounterOut, err := db_trustlines.GetSyncCounterOut(session.Datagram)
    if err != nil {
        log.Printf("Error getting sync_counter_out: %v", err)
        return
    }

    dg := main.Datagram{
        Command:           main.Server_SetTrustlineSyncTimestamp,
        Username:          session.Datagram.PeerUsername,
        PeerUsername:      session.Datagram.Username,
        PeerServerAddress: main.GetServerAddress(),
    }
    binary.BigEndian.PutUint32(dg.Counter[:], syncCounterOut)

    if err := handlers.SignAndSendDatagram(session, &dg); err != nil {
        log.Printf("Failed to sign and send datagram: %v", err)
        return
    }
    log.Println("SetTrustlineSyncTimestamp command sent successfully.")
}

func sendTrustline(session main.Session, counter uint32) {
    trustline, err := db_trustlines.GetTrustlineOut(session.Datagram)
    if err != nil {
        log.Printf("Error getting trustline: %v", err)
        return
    }

    dg := main.Datagram{
        Command:           main.Server_SetTrustline,
        Username:          session.Datagram.PeerUsername,
        PeerUsername:      session.Datagram.Username,
        PeerServerAddress: main.GetServerAddress(),
    }
    binary.BigEndian.PutUint32(dg.Arguments[:4], trustline)
    binary.BigEndian.PutUint32(dg.Counter[:], counter)

    if err := handlers.SignAndSendDatagram(session, &dg); err != nil {
        log.Printf("Failed to sign and send datagram: %v", err)
        return
    }
    log.Println("SetTrustline command sent successfully.")
}
