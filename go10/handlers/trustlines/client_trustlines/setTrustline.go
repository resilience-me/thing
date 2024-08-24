package client_trustlines

import (
    "encoding/binary"
    "log"

    "ripple/comm"
    "ripple/database/db_trustlines"
    "ripple/types"
    "ripple/handlers/trustlines"
)

// SetTrustline updates the trustline based on the given session.
func SetTrustline(session types.Session) {
    datagram := session.Datagram

    // Retrieve the trustline amount from the Datagram
    trustlineAmount := binary.BigEndian.Uint32(datagram.Arguments[:4])

    // Write the new trustline amount using the setter in db_trustlines
    if err := db_trustlines.GetTrustlineInFromDatagram(datagram, trustlineAmount); err != nil {
        log.Printf("Error writing trustline to file for user %s: %v", datagram.Username, err)
        comm.SendErrorResponse(session.Addr, "Failed to write trustline.")
        return
    }

    // Increment the sync_counter using the function in trustlines package
    if err := trustlines.IncrementSyncCounter(datagram); err != nil {
        log.Printf("Error incrementing sync_counter for user %s: %v", datagram.Username, err)
        comm.SendErrorResponse(session.Addr, "Failed to update sync counter.")
        return
    }

    // Log success
    log.Printf("Trustline and sync counter updated successfully for user %s.", datagram.Username)

    // Send success response
    if err := comm.SendSuccessResponse(session.Addr, []byte("Trustline updated successfully.")); err != nil {
        log.Printf("Failed to send success response to user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("Sent success response to client for user %s.", datagram.Username)
}
