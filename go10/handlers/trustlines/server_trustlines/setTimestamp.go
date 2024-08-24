package server_trustlines

import (
    "log"
    "time"
    "ripple/main"
    "ripple/database/db_trustlines"
)

// SetTimestamp handles updating the sync timestamp for trustlines
func SetTimestamp(session main.Session) {
    datagram := session.Datagram

    // Retrieve the current timestamp
    timestamp := time.Now().Unix()

    // Write the new timestamp using the setter in db_trustlines
    if err := db_trustlines.SetTimestamp(datagram, timestamp); err != nil {
        log.Printf("Error writing timestamp for user %s: %v", datagram.Username, err)
        return
    }

    // Log success
    log.Printf("Timestamp and counter_in updated successfully for user %s.", datagram.Username)
}
