package comm

import (
	"fmt"
	"ripple/udpr"
)

// SendSuccessResponse sends a success message using the provided Conn with retry logic.
func SendSuccessResponse(data []byte, conn *Conn) error {
	response := append([]byte{0}, data...) // Combine success indicator and message

	if err := udpr.SendWithRetry(conn.UDPConn, conn.addr, response, HighImportance); err != nil {
		return fmt.Errorf("error sending success response: %w", err) // Return detailed error
	}

	return nil
}

// SendErrorResponse sends an error message using the provided Conn with retry logic.
func SendErrorResponse(message string, conn *Conn) error {
	response := append([]byte{1}, []byte(message)...) // Combine error indicator and message

	if err := udpr.SendWithRetry(conn.UDPConn, conn.addr, response, HighImportance); err != nil {
		return fmt.Errorf("error sending error response: %w", err) // Return detailed error
	}

	return nil
}
