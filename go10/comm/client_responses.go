package comm

import "fmt"

// SendSuccessResponse sends a success message using the provided address with retry logic.
func SendSuccessResponse(data []byte, addr *net.UDPAddr) error {
	response := append([]byte{0}, data...) // Combine success indicator and message
	if err := SendWithAddress(addr, response, HighImportance); err != nil {
		return fmt.Errorf("error sending success response: %w", err)
	}
	return nil
}

// SendErrorResponse sends an error message using the provided address with retry logic.
func SendErrorResponse(message string, addr *net.UDPAddr) error {
	response := append([]byte{1}, []byte(message)...) // Combine error indicator and message
	if err := SendWithAddress(addr, response, HighImportance); err != nil {
		return fmt.Errorf("error sending error response: %w", err)
	}
	return nil
}
