package client_trustlines

import (
    "fmt"
    "ripple/database/db_trustlines" // Updated to match your import structure
    "ripple/main"                   // Updated to match your import structure
)

// GetTrustlineIn handles fetching the inbound trustline information
func GetTrustlineIn(session main.Session) {
    // Fetch the inbound trustline
    trustline, err := db_trustlines.GetTrustlineIn(session.Datagram)
    if err != nil {
        fmt.Printf("Error reading inbound trustline: %v\n", err) // Log the error
        main.SendErrorResponse("Error reading inbound trustline.", session.Conn)
        return
    }

    // Prepare success response using the renamed function
    responseData := uint32ToBytes(trustline)

    // Send the success response back to the client
    if err := main.SendSuccessResponse(responseData, session.Conn); err != nil {
        fmt.Printf("Error sending success response: %v\n", err) // Log the error
        return
    }

    fmt.Printf("Inbound trustline sent successfully.\n")
}
