package comm

import (
	"fmt"
	"net"
)

// SendSuccessResponse sends a success message using the provided address with retry logic.
func SendSuccessResponse(addr *net.UDPAddr, data []byte) error {
	response := append([]byte{0}, data...) // Combine success indicator and message
	if err := SendWithAddress(addr, response, HighImportance); err != nil {
		return fmt.Errorf("error sending success response: %w", err)
	}
	return nil
}

// SendErrorResponse sends an error message using the provided address with retry logic.
func SendErrorResponse(addr *net.UDPAddr, message string) error {
	response := append([]byte{1}, []byte(message)...) // Combine error indicator and message
	if err := SendWithAddress(addr, response, HighImportance); err != nil {
		return fmt.Errorf("error sending error response: %w", err)
	}
	return nil
}
