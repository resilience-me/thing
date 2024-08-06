package client_trustlines

import (
    "encoding/binary"
    "fmt"
    "resilience/database/db_trustlines"
    "resilience/handlers"
    "resilience/main"
)

// GetTrustlineIn handles fetching the inbound trustline information
func GetTrustlineIn(ctx main.HandlerContext) {
    // Validate the client request
    if errorMessage, err := handlers.ValidateClientRequest(ctx.Datagram); err != nil {
        fmt.Printf("Validation failed: %v\n", err) // Log detailed error
        _ = handlers.SendErrorResponse(ctx, errorMessage)
        return
    }

    // Fetch the inbound trustline
    trustline, err := db_trustlines.GetTrustlineIn(ctx.Datagram)
    if err != nil {
        fmt.Printf("Error reading inbound trustline: %v\n", err) // Log the error
        _ = handlers.SendErrorResponse(ctx, "Error reading inbound trustline.")
        return
    }

    // Prepare success response
    responseData := make([]byte, 4) // Allocate 4 bytes for the trustline
    binary.BigEndian.PutUint32(responseData, trustline) // Convert the trustline to bytes

    // Send the success response back to the client
    if err := handlers.SendSuccessResponse(ctx, responseData); err != nil {
        fmt.Printf("Error sending success response: %v\n", err) // Log the error
        return
    }

    fmt.Printf("Inbound trustline sent successfully.\n")
}
