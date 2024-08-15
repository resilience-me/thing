package server_trustlines

import (
    "log"
    "time"
    "ripple/main"
    "ripple/trustlines"             // Import the trustlines package for counter validation
    "ripple/database/db_trustlines" // Handles database-related operations
    "ripple/database/db_server" // Handles database-related operations
)

// SetTimestamp handles updating the sync timestamp for trustlines
func SetTimestamp(session main.Session) {
    datagram := session.Datagram

    // Validate the counter_in using the ValidateCounterIn function from trustlines package
    if err := db_server.ValidateCounterIn(datagram); err != nil {
        log.Printf("Counter_in validation failed for user %s: %v", datagram.Username, err)
        return
    }

    // Retrieve the current timestamp
    timestamp := time.Now().Unix()

    // Write the new timestamp using the setter in db_trustlines
    if err := db_trustlines.SetTimestamp(datagram, timestamp); err != nil {
        log.Printf("Error writing timestamp for user %s: %v", datagram.Username, err)
        return
    }

    // After successfully updating the timestamp, update the counter_in
    if err := db_server.SetCounterIn(datagram, datagram.Counter); err != nil {
        log.Printf("Error updating counter_in for user %s: %v", datagram.Username, err)
        return
    }

    // Log success
    log.Printf("Timestamp and counter_in updated successfully for user %s.", datagram.Username)
}
