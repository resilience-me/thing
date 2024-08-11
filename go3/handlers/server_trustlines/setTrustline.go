package server_trustlines

import (
    "log"
    "time"
    "ripple/main"
    "ripple/handlers"
    "ripple/database/db_trustlines"
)

// SetTrustline handles setting or updating a trustline from another server's perspective
func SetTrustline(session main.Session) {
    // Retrieve the counter_in value using the new getter
    counterIn, err := db_trustlines.GetCounterIn(session.Datagram)
    if err != nil {
        log.Printf("Error getting counter_in for user %s: %v", session.Datagram.Username, err)
        return
    }

    // Check the counter
    counter := session.Datagram.Counter
    if counter <= counterIn {
        log.Printf("Potential replay attack: received counter (%d) is not greater than counter_in (%d) for user %s.", counter, counterIn, session.Datagram.Username)
        return
    }

    // Retrieve the sync_in value using the new getter
    prevSyncIn, err := db_trustlines.GetSyncIn(session.Datagram)
    if err != nil {
        log.Printf("Error getting sync_in for user %s: %v", session.Datagram.Username, err)
        return
    }

    // Retrieve the syncIn counter from the Datagram
    syncInBytes := session.Datagram.Arguments[4:8]
    syncIn := main.BytesToUint32(syncInBytes)

    if syncIn > prevSyncIn {
        handleTrustlineUpdate(session, syncInBytes, syncIn)
    } else {
        handleTimestampOnly(session)
    }
}

func handleTrustlineUpdate(session main.Session, syncInBytes []byte, syncIn uint32) {
    // Retrieve the trustline amount from the Datagram
    trustlineAmount := main.BytesToUint32(session.Datagram.Arguments[:4])

    if err := db_trustlines.SetTrustlineIn(session.Datagram, trustlineAmount); err != nil {
        log.Printf("Error writing trustline to file for user %s: %v", session.Datagram.Username, err)
        return
    }

    if err := db_trustlines.SetSyncIn(session.Datagram, syncIn); err != nil {
        log.Printf("Error writing sync_in to file for user %s: %v", session.Datagram.Username, err)
        return
    }

    if err := db_trustlines.SetTimestamp(session.Datagram, time.Now().Unix()); err != nil {
        log.Printf("Error writing timestamp to file for user %s: %v", session.Datagram.Username, err)
        return
    }
    log.Printf("Trustline, sync_in, and timestamp updated successfully for user %s.", session.Datagram.Username)

    // Retrieve and increment the counter_out value
    counterOut, err := db_trustlines.GetAndIncrementCounterOut(session.Datagram)
    if err != nil {
        log.Printf("Error handling counter_out for user %s: %v", session.Datagram.Username, err)
        return
    }

    // Prepare the datagram to send back to the peer
    dg := main.Datagram{
        Command:           main.ServerTrustlines_SetSyncOut,
        Username:          session.Datagram.PeerUsername,
        PeerUsername:      session.Datagram.Username,
        PeerServerAddress: main.GetServerAddress(),
        Arguments:         syncInBytes,
        Counter:           counterOut,
    }

    if err := handlers.SignAndSendDatagram(session, &dg); err != nil {
        log.Printf("Failed to sign and send datagram for user %s: %v", session.Datagram.Username, err)
        return
    }

    log.Printf("Trustline update and datagram sent successfully for user %s.", session.Datagram.Username)
}

func handleTimestampOnly(session main.Session) {
    log.Printf("Sync_in is synchronized with the peer's most recent trustline_out for user %s.", session.Datagram.Username)

    if err := db_trustlines.SetTimestamp(session.Datagram, time.Now().Unix()); err != nil {
        log.Printf("Error writing timestamp to file for user %s: %v", session.Datagram.Username, err)
        return
    }

    log.Printf("Trustline synchronization timestamp updated successfully for user %s.", session.Datagram.Username)
}
