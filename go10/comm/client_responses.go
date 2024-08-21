package comm

import (
	"fmt"
	"ripple/udpr"
)

// SendSuccessResponse sends a success message using the provided address with retry logic.
func SendSuccessResponse(data []byte, addr *net.UDPAddr) error {
	response := append([]byte{0}, data...) // Combine success indicator and message

	// Create a UDP connection with an ephemeral local port
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %w", err)
	}
	defer conn.Close()

	if err := udpr.SendWithRetry(conn, addr, response, HighImportance); err != nil {
		return fmt.Errorf("error sending success response: %w", err) // Return detailed error
	}

	return nil
}

// SendErrorResponse sends an error message using the provided address with retry logic.
func SendErrorResponse(message string, addr *net.UDPAddr) error {
	response := append([]byte{1}, []byte(message)...) // Combine error indicator and message

	// Create a UDP connection with an ephemeral local port
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to create UDP connection: %w", err)
	}
	defer conn.Close()

	if err := udpr.SendWithRetry(conn, addr, response, HighImportance); err != nil {
		return fmt.Errorf("error sending error response: %w", err) // Return detailed error
	}

	return nil
}
