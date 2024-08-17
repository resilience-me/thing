package main

import (
    "fmt"
    "net"
)

// SendSuccessResponse sends a success message to the specified UDP address.
func SendSuccessResponse(data []byte, conn *net.UDPConn, addr *net.UDPAddr) error {
    response := append([]byte{0}, data...) // Combine success indicator and message

    if _, err := conn.WriteToUDP(response, addr); err != nil { // Send combined response to the UDP address
        return fmt.Errorf("error sending success response: %w", err) // Return detailed error
    }

    return nil
}

// SendErrorResponse sends an error message to the specified UDP address.
func SendErrorResponse(message string, conn *net.UDPConn, addr *net.UDPAddr) error {
    response := append([]byte{1}, []byte(message)...) // Combine error indicator (e.g., '1') and message

    if _, err := conn.WriteToUDP(response, addr); err != nil { // Send combined response to the UDP address
        return fmt.Errorf("error sending error response: %w", err) // Return detailed error
    }

    return nil
}
