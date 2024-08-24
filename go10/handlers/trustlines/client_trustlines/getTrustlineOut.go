package client_trustlines

import (
    "log"

    "ripple/comm"
    "ripple/database/db_trustlines"
    "ripple/types"
)

// GetTrustlineOut handles fetching the outbound trustline information
func GetTrustlineOut(session types.Session) {
    datagram := session.Datagram

    // Fetch the outbound trustline
    trustline, err := db_trustlines.GetTrustlineOutFromDatagram(datagram)
    if err != nil {
        log.Printf("Error reading outbound trustline for user %s: %v", datagram.Username, err)
        comm.SendErrorResponse(session.Addr, "Error reading outbound trustline.")
        return
    }

    // Prepare success response using the main utility function
    responseData := types.Uint32ToBytes(trustline)

    // Send the success response back to the client
    if err := comm.SendSuccessResponse(session.Addr, responseData); err != nil {
        log.Printf("Error sending success response to user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("Outbound trustline sent successfully to user %s.", datagram.Username) // Log successful operation
}
