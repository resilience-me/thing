package client_trustlines

import (
    "encoding/binary"
    "fmt"
    "net"
    "os"
    "path/filepath"
    "strconv"

    "resilience/main"
    "resilience/handlers"
	"resilience/database/db_trustlines"
)

// getTrustline handles fetching the trustline information for both inbound and outbound.
func getTrustline(ctx main.HandlerContext, trustline uint32) {
    // Validate the client request
    if err := handlers.ValidateClientRequest(ctx); err != nil {
        fmt.Printf("Validation failed: %v\n", err) // Log detailed error
        return // Error response has already been sent in ValidateClientRequest
    }

    // Prepare success response
    responseData := make([]byte, 4) // Allocate 4 bytes for the trustline amount
    binary.BigEndian.PutUint32(responseData, uint32(trustline)) // Convert the trustline amount to bytes

    // Send the success response back to the client
    if err := handlers.SendSuccessResponse(ctx, responseData); err != nil {
        fmt.Printf("Error sending success response: %v\n", err) // Log the error
        return
    }

    fmt.Printf("Trustline amount (%s) sent successfully.\n", filename)
}

// GetTrustlineIn handles fetching the inbound trustline information
func GetTrustlineIn(ctx main.HandlerContext) {
    trustline := db_trustlines.GetTrustlineIn(ctx.Datagram)
    if err != nil {
        fmt.Printf("Error reading inbound trustline amount: %v\n", err) // Log the error
        _ = handlers.SendErrorResponse(ctx, "Error reading inbound trustline.")
        return
    }
    getTrustline(ctx, trustline)
}

// GetTrustlineOut handles fetching the outbound trustline information
func GetTrustlineOut(ctx main.HandlerContext) {
    trustline := db_trustlines.GetTrustlineOut(ctx.Datagram)
    if err != nil {
        fmt.Printf("Error reading outbound trustline amount: %v\n", err) // Log the error
        _ = handlers.SendErrorResponse(ctx, "Error reading outbound trustline.")
        return
    }
    getTrustline(ctx, trustline)
}
