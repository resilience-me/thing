package client_trustlines

import (
    "log"
    "ripple/database/db_trustlines" // Updated to match your import structure
    "ripple/main"                   // Updated to match your import structure
)

// GetTrustlineIn handles fetching the inbound trustline information
func GetTrustlineIn(session main.Session) {
    datagram := session.Datagram

    // Retrieve the previous client-side counter value using the getter
    prevCounter, err := db_trustlines.GetCounter(datagram)
    if err != nil {
        log.Printf("Error getting previous counter for user %s: %v", datagram.Username, err)
        main.SendErrorResponse("Failed to read counter file.", session.Conn)
        return
    }

    // Check if the client-side counter is valid (prevents replay attacks)
    if datagram.Counter <= prevCounter {
        log.Printf("Received counter is not greater than previous counter for user %s. Potential replay attack.", datagram.Username)
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
    responseData := uint32ToBytes(trustline)

    // Send the success response back to the client
    if err := main.SendSuccessResponse(responseData, session.Conn); err != nil {
        log.Printf("Error sending success response to user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("Inbound trustline sent successfully to user %s.", datagram.Username)
}
