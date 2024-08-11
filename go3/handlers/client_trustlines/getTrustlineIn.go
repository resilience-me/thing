package client_trustlines

import (
    "log"
    "ripple/database/db_trustlines" // Updated to match your import structure
    "ripple/main"                   // Updated to match your import structure
)

// GetTrustlineIn handles fetching the inbound trustline information
func GetTrustlineIn(session main.Session) {
    // Fetch the inbound trustline
    trustline, err := db_trustlines.GetTrustlineIn(session.Datagram)
    if err != nil {
        log.Printf("Error reading inbound trustline for user %s: %v", session.Datagram.Username, err) // Log the error with context
        main.SendErrorResponse("Error reading inbound trustline.", session.Conn)
        return
    }

    // Prepare success response using the renamed function
    responseData := uint32ToBytes(trustline)

    // Send the success response back to the client
    if err := main.SendSuccessResponse(responseData, session.Conn); err != nil {
        log.Printf("Error sending success response to user %s: %v", session.Datagram.Username, err) // Log the error with context
        return
    }

    log.Printf("Inbound trustline sent successfully to user %s.", session.Datagram.Username) // Log successful operation
}
