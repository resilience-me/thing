package client_trustlines

import (
    "log"
    "ripple/database/db_trustlines" // Handles database-related operations
    "ripple/main"                   // Main package for session and communication utilities
    "ripple/trustlines"             // Import the trustlines package for counter validation
)

// GetTrustlineIn handles fetching the inbound trustline information
func GetTrustlineIn(session main.Session) {
    datagram := session.Datagram

    // Validate the counter using the ValidateCounter function from trustlines package
    if err := trustlines.ValidateCounter(datagram); err != nil {
        log.Printf("Counter validation failed for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Received counter is not valid.", session.Conn)
        return
    }

    // Fetch the inbound trustline
    trustline, err := db_trustlines.GetTrustlineIn(datagram)
    if err != nil {
        log.Printf("Error reading inbound trustline for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Error reading inbound trustline.", session.Conn)
        return
    }

    // Update the counter value after validation
    if err := db_trustlines.SetCounter(datagram, datagram.Counter); err != nil {
        log.Printf("Error updating counter for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to update counter.", session.Conn)
        return
    }

    // Prepare success response
    responseData := main.Uint32ToBytes(trustline)

    // Send the success response back to the client
    if err := main.SendSuccessResponse(responseData, session.Conn); err != nil {
        log.Printf("Error sending success response to user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("Inbound trustline sent successfully to user %s.", datagram.Username)
}
