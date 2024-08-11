package client_trustlines

import (
    "log"
    "ripple/database/db_trustlines" // Updated to match your import structure
    "ripple/main"                   // Updated to match your import structure
)

// GetTrustlineOut handles fetching the outbound trustline information
func GetTrustlineOut(session main.Session) {
    // Fetch the outbound trustline
    trustline, err := db_trustlines.GetTrustlineOut(session.Datagram)
    if err != nil {
        log.Printf("Error reading outbound trustline for user %s: %v", session.Datagram.Username, err) // Log the error with context
        _ = main.SendErrorResponse([]byte("Error reading outbound trustline."), session.Conn)
        return
    }

    // Prepare success response using the renamed function
    responseData := uint32ToBytes(trustline)

    // Send the success response back to the client
    if err := main.SendSuccessResponse(responseData, session.Conn); err != nil {
        log.Printf("Error sending success response to user %s: %v", session.Datagram.Username, err) // Log the error with context
        return
    }

    log.Printf("Outbound trustline sent successfully to user %s.", session.Datagram.Username) // Log successful operation
}
