package main

import (
    "fmt"
    "net"
)

// SendSuccessResponse sends a success message using the provided Conn.
func SendSuccessResponse(data []byte, conn *Conn) error {
    response := append([]byte{0}, data...) // Combine success indicator and message

    if _, err := conn.conn.WriteToUDP(response, conn.addr); err != nil { // Send combined response using the stored UDP connection and address
        return fmt.Errorf("error sending success response: %w", err) // Return detailed error
    }

    return nil
}

// SendErrorResponse sends an error message using the provided Conn.
func SendErrorResponse(message string, conn *Conn) error {
    response := append([]byte{1}, []byte(message)...) // Combine error indicator (e.g., '1') and message

    if _, err := conn.conn.WriteToUDP(response, conn.addr); err != nil { // Send combined response using the stored UDP connection and address
        return fmt.Errorf("error sending error response: %w", err) // Return detailed error
    }

    return nil
}
