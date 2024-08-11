package server_trustlines

import (
    "encoding/binary"
    "log"
    "time"
    "ripple/main"
    "ripple/handlers"
    "ripple/database/db_trustlines"
)

// SetTimestamp handles updating the sync timestamp for trustlines
func SetTimestamp(session main.Session) {
    // Retrieve the previous counter_in value
    prevCounterIn, err := db_trustlines.GetCounterIn(session.Datagram)
    if err != nil {
        log.Printf("Error getting previous counter_in for user %s: %v", session.Datagram.Username, err)
        return
    }

    // Get the new counter value from the datagram
    counter := binary.BigEndian.Uint32(session.Datagram.Counter[:])

    // Check if the new counter is greater than the previous counter_in
    if counter <= prevCounterIn {
        log.Printf("Received counter (%d) is not greater than previous counter_in (%d) for user %s. Potential replay attack.",
            counter, prevCounterIn, session.Datagram.Username)
        return
    }

    // Write the new counter_in value
    if err := db_trustlines.SetCounterIn(session.Datagram, counter); err != nil {
        log.Printf("Error writing counter_in for user %s: %v", session.Datagram.Username, err)
        return
    }

    // Retrieve the current timestamp
    timestamp := time.Now().Unix()

    // Write the new timestamp
    if err := db_trustlines.SetTimestamp(session.Datagram, timestamp); err != nil {
        log.Printf("Error writing timestamp for user %s: %v", session.Datagram.Username, err)
        return
    }

    log.Printf("counter_in and timestamp updated successfully for user %s.", session.Datagram.Username)
}
