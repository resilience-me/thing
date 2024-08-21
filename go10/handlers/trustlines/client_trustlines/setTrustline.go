package client_trustlines

import (
    "encoding/binary"
    "log"
    "ripple/database/db_trustlines"
    "ripple/main"
    "ripple/comm"
    "ripple/trustlines"
)

// SetTrustline updates the trustline based on the given session.
func SetTrustline(session main.Session) {
    datagram := session.Datagram

    // Retrieve the trustline amount from the Datagram
    trustlineAmount := binary.BigEndian.Uint32(datagram.Arguments[:4])

    // Write the new trustline amount using the setter in db_trustlines
    if err := db_trustlines.SetTrustlineOut(datagram, trustlineAmount); err != nil {
        log.Printf("Error writing trustline to file for user %s: %v", datagram.Username, err)
        comm.SendErrorResponse("Failed to write trustline.", session.Conn)
        return
    }

    // Increment the sync_counter using the function in trustlines package
    if err := trustlines.IncrementSyncCounter(datagram); err != nil {
        log.Printf("Error incrementing sync_counter for user %s: %v", datagram.Username, err)
        comm.SendErrorResponse("Failed to update sync counter.", session.Conn)
        return
    }

    // Log success
    log.Printf("Trustline and sync counter updated successfully for user %s.", datagram.Username)

    // Send success response
    if err := comm.SendSuccessResponse([]byte("Trustline updated successfully."), session.Conn); err != nil {
        log.Printf("Failed to send success response to user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("Sent success response to client for user %s.", datagram.Username)
}
