package client_trustlines

import (
    "log"

    "ripple/comm"
    "ripple/database/db_trustlines"
    "ripple/main"
    "ripple/types"
)

// GetTrustlineOut handles fetching the outbound trustline information
func GetTrustlineOut(session main.Session) {
    datagram := session.Datagram

    // Fetch the outbound trustline
    trustline, err := db_trustlines.GetTrustlineOut(datagram)
    if err != nil {
        log.Printf("Error reading outbound trustline for user %s: %v", datagram.Username, err)
        comm.SendErrorResponse("Error reading outbound trustline.", session.Conn)
        return
    }

    // Prepare success response using the main utility function
    responseData := types.Uint32ToBytes(trustline)

    // Send the success response back to the client
    if err := comm.SendSuccessResponse(responseData, session.Addr); err != nil {
        log.Printf("Error sending success response to user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("Outbound trustline sent successfully to user %s.", datagram.Username) // Log successful operation
}
