package handlers

import (
	"fmt"
	"ripple/udpr"
)

// SendSuccessResponse sends a success message using the provided Conn with retry logic.
func SendSuccessResponse(client *Client, data []byte) error {
	response := append([]byte{0}, data...) // Combine success indicator and message

	// Use the SendClient wrapper for standard priority
	if err := udpr.SendClient(client, response); err != nil {
		return fmt.Errorf("error sending success response: %w", err) // Return detailed error
	}

	return nil
}

// SendErrorResponse sends an error message using the provided Conn with retry logic.
func SendErrorResponse(client *Client, message string) error {
	response := append([]byte{1}, []byte(message)...) // Combine error indicator and message

	// Use the SendClient wrapper for standard priority
	if err := udpr.SendClient(client, response); err != nil {
		return fmt.Errorf("error sending error response: %w", err) // Return detailed error
	}

	return nil
}
