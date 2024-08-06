package client_trustlines

import (
    "encoding/binary"
    "fmt"
    "resilience/database/db_trustlines"
    "resilience/handlers"
    "resilience/main"
)

// GetTrustlineOut handles fetching the outbound trustline information
func GetTrustlineOut(ctx main.HandlerContext) {
    // Validate the client request
    if err := handlers.ValidateClientRequest(ctx); err != nil {
        fmt.Printf("Validation failed: %v\n", err) // Log detailed error
        return // Error response has already been sent in ValidateClientRequest
    }

    // Fetch the outbound trustline
    trustline, err := db_trustlines.GetTrustlineOut(ctx.Datagram)
    if err != nil {
        fmt.Printf("Error reading outbound trustline: %v\n", err) // Log the error
        _ = handlers.SendErrorResponse(ctx, "Error reading outbound trustline.")
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

    fmt.Printf("Outbound trustline sent successfully.\n")
}
