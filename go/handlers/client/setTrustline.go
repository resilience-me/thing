package client

import (
	"encoding/binary"
	"fmt"

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

	// Retrieve the previous counter value using the getter
	prevCounter, err := main.GetCounter(ctx.Datagram)
	if err != nil {
		fmt.Printf("Error getting previous counter: %v\n", err) // Log detailed error
		_ = handlers.SendErrorResponse(ctx, "Failed to read counter file.") // Send simpler error message
		return
	}

	// Check the counter
	counter := binary.BigEndian.Uint32(ctx.Datagram.Counter[:])
	if counter <= prevCounter {
		fmt.Println("Received counter is not greater than previous counter. Potential replay attack.")
		_ = handlers.SendErrorResponse(ctx, "Received counter is not valid.") // Send simpler error message
		return
	}

	// Retrieve the trustline amount from the Datagram
	trustlineAmount := binary.BigEndian.Uint32(ctx.Datagram.Arguments[:4])

	// Write the new trustline amount using the setter
	if err := main.SetTrustlineOut(ctx.Datagram, trustlineAmount); err != nil {
		fmt.Printf("Error writing trustline to file: %v\n", err) // Log detailed error
		_ = handlers.SendErrorResponse(ctx, "Failed to write trustline.") // Send simpler error message
		return
	}

	// Write the new counter value using the setter
	if err := main.SetCounter(ctx.Datagram, counter); err != nil {
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
