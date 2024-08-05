package client

import (
    "encoding/binary"
    "fmt"
    "net"
    "os"
    "path/filepath"
    "strconv"

    "resilience/main"
    "resilience/handlers" // Import the handlers package
)

// SetTrustline handles setting or updating a trustline from the client's perspective
func SetTrustline(ctx main.HandlerContext) {
    // Validate the client request (account and peer directory checks, and signature verification)
    if err := handlers.ValidateClientRequest(ctx); err != nil {
        fmt.Printf("Validation failed: %v\n", err) // Log detailed error
        return // Error response has already been sent in ValidateClientRequest
    }

    // Get the trustline directory
    trustlineDir := main.GetTrustlineDir(ctx.Datagram)

    // Construct the trustline and counter file paths
    counterOutPath := filepath.Join(trustlineDir, "counter_out.txt")
    trustlineOutPath := filepath.Join(trustlineDir, "trustline_out.txt")

    // Load the previous counter value
    prevCounterStr, err := os.ReadFile(counterOutPath)
    if err != nil && !os.IsNotExist(err) {
        fmt.Printf("Error reading counter file: %v\n", err) // Log detailed error
        _ = handlers.SendErrorResponse(ctx, "Failed to read counter file.") // Send simpler error message
        return
    }

    // Parse previous counter
    prevCounter, err := strconv.ParseUint(string(prevCounterStr), 10, 32) // Parse as uint64 first
    if err != nil {
        fmt.Printf("Error parsing previous counter string: %v\n", err) // Log detailed error
        _ = handlers.SendErrorResponse(ctx, "Failed to parse previous counter.") // Send simpler error message
        return
    }

    // Check the counter
    counter := binary.BigEndian.Uint32(ctx.Datagram.Counter[:])
    if counter <= uint32(prevCounter) {
        fmt.Println("Received counter is not greater than previous counter. Potential replay attack.")
        _ = handlers.SendErrorResponse(ctx, "Received counter is not valid.") // Send simpler error message
        return
    }

    trustlineAmount := binary.BigEndian.Uint32(ctx.Datagram.Arguments[:4])
    
    // Write the new trustline amount to the file
    if err := os.WriteFile(trustlineOutPath, []byte(fmt.Sprintf("%d", trustlineAmount)), 0644); err != nil {
        fmt.Printf("Error writing trustline to file: %v\n", err) // Log detailed error
        _ = handlers.SendErrorResponse(ctx, "Failed to write trustline.") // Send simpler error message
        return
    }

    // Write the new counter value as a string
    counterStr := fmt.Sprintf("%d", counter)
    if err := os.WriteFile(counterOutPath, []byte(counterStr), 0644); err != nil {
        fmt.Printf("Error writing counter to file: %v\n", err) // Log detailed error
        _ = handlers.SendErrorResponse(ctx, "Failed to write counter.") // Send simpler error message
        return
    }

    fmt.Println("Trustline and counter updated successfully.")

    // Prepare success response
    successMessage := []byte("Trustline updated successfully.")
    if err := handlers.SendSuccessResponse(ctx, successMessage); err != nil {
        fmt.Printf("Error sending success response: %v\n", err) // Log detailed error
        return
    }
    fmt.Println("Sent success response to client.")
}
