package client_trustlines

import (
    "encoding/binary"
    "log"
    "ripple/database/db_trustlines"
    "ripple/main"
    "ripple/trustlines" // Import the trustlines package for counter validation and incrementing sync counter
)

// SetTrustline updates the trustline based on the given session.
func SetTrustline(session main.Session) {
    datagram := session.Datagram

    // Validate the counter using the ValidateCounter function from trustlines package
    if err := trustlines.ValidateCounter(datagram); err != nil {
        log.Printf("Counter validation failed for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Received counter is not valid.", session.Conn)
        return
    }

    // Retrieve the trustline amount from the Datagram
    trustlineAmount := binary.BigEndian.Uint32(datagram.Arguments[:4])

    // Write the new trustline amount using the setter in db_trustlines
    if err := db_trustlines.SetTrustlineOut(datagram, trustlineAmount); err != nil {
        log.Printf("Error writing trustline to file for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to write trustline.", session.Conn)
        return
    }

    // Increment the sync_counter using the function in trustlines package
    if err := trustlines.IncrementSyncCounter(datagram); err != nil {
        log.Printf("Error incrementing sync_counter for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to update sync counter.", session.Conn)
        return
    }

    // Write the new client-side counter value using the setter in db_trustlines
    if err := db_trustlines.SetCounter(datagram, datagram.Counter); err != nil {
        log.Printf("Error writing counter to file for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to write counter.", session.Conn)
        return
    }

    // Log success
    log.Printf("Trustline and counters updated successfully for user %s.", datagram.Username)

    // Send success response
    if err := main.SendSuccessResponse([]byte("Trustline updated successfully."), session.Conn); err != nil {
        log.Printf("Failed to send success response to user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("Sent success response to client for user %s.", datagram.Username)
}
