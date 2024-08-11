package main

import (
    "fmt"
    "net"
)

// SendSuccessResponse sends a success message over the given connection.
func SendSuccessResponse(message string, conn net.Conn) error {
    response := append([]byte{0}, []byte(message)...) // Combine success indicator and message

    if _, err := conn.Write(response); err != nil { // Send combined response
        return fmt.Errorf("error sending success response: %w", err) // Return detailed error
    }

    return nil
}

// SendErrorResponse sends an error message over the given connection.
func SendErrorResponse(message string, conn net.Conn) error {
    response := append([]byte{1}, []byte(message)...) // Combine error indicator (e.g., '1') and message

    if _, err := conn.Write(response); err != nil { // Send combined response
        return fmt.Errorf("error sending error response: %w", err) // Return detailed error
    }

    return nil
}
