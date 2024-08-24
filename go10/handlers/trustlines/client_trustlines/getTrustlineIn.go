package client_trustlines

import (
    "log"

    "ripple/comm"
    "ripple/database/db_trustlines"
    "ripple/types"
)

// GetTrustlineIn handles fetching the inbound trustline information
func GetTrustlineIn(session types.Session) {
    datagram := session.Datagram

    // Fetch the inbound trustline
    trustline, err := db_trustlines.GetTrustlineInFromDatagram(datagram)
    if err != nil {
        log.Printf("Error reading inbound trustline for user %s: %v", datagram.Username, err)
        comm.SendErrorResponse(session.Addr, "Error reading inbound trustline.")
        return
    }

    // Prepare success response
    responseData := types.Uint32ToBytes(trustline)

    // Send the success response back to the client
    if err := comm.SendSuccessResponse(session.Addr, responseData); err != nil {
        log.Printf("Error sending success response to user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("Inbound trustline sent successfully to user %s.", datagram.Username)
}
