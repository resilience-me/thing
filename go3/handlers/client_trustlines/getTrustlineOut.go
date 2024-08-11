package client_trustlines

import (
    "fmt"
    "ripple/database/db_trustlines" // Updated to match your import structure
    "ripple/main"                   // Updated to match your import structure
)

// GetTrustlineOut handles fetching the outbound trustline information
func GetTrustlineOut(session main.Session) {
    // Fetch the outbound trustline
    trustline, err := db_trustlines.GetTrustlineOut(session.Datagram)
    if err != nil {
        fmt.Printf("Error reading outbound trustline: %v\n", err) // Log the error
        _ = main.SendErrorResponse([]byte("Error reading outbound trustline."), session.Conn)
        return
    }

    // Prepare success response using the renamed function
    responseData := uint32ToBytes(trustline)

    // Send the success response back to the client
    if err := main.SendSuccessResponse(responseData, session.Conn); err != nil {
        fmt.Printf("Error sending success response: %v\n", err) // Log the error
        return
    }

    fmt.Printf("Outbound trustline sent successfully.\n")
}
