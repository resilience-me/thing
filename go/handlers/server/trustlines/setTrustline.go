package trustlines

import (
	"encoding/binary"
	"fmt"
	"time"
	"resilience/main"
	"resilience/handlers"
)

// SetTrustline handles setting or updating a trustline from another server's perspective
func SetTrustline(ctx main.HandlerContext) {

	if err := handlers.ValidateServerRequest(ctx); err != nil {
		fmt.Printf("Validation failed: %v\n", err) // Log detailed error
		return
	}

	// Retrieve the sync_in value using the new getter
	syncIn, err := main.GetSyncIn(ctx.Datagram)
	if err != nil {
		fmt.Printf("Error getting sync_in: %v\n", err)
		return
	}

	// Check the counter
	counter := binary.BigEndian.Uint32(ctx.Datagram.Counter[:])
	if counter <= syncIn {
		fmt.Println("Received counter is not greater than sync_in. Potential replay attack.")
		return
	}

	// Retrieve the trustline amount from the Datagram
	trustlineAmount := binary.BigEndian.Uint32(ctx.Datagram.Arguments[:4])

	// Write the new trustline amount using the setter
	if err := main.SetTrustlineOut(ctx.Datagram, trustlineAmount); err != nil {
		fmt.Printf("Error writing trustline to file: %v\n", err)
		return
	}

	// Write the new sync_in value using the setter
	if err := main.SetSyncIn(ctx.Datagram, counter); err != nil {
		fmt.Printf("Error writing sync_in to file: %v\n", err)
		return
	}

	// Write the Unix timestamp using the setter
	if err := main.SetTimestamp(ctx.Datagram, time.Now().Unix()); err != nil {
		fmt.Printf("Error writing timestamp to file: %v\n", err)
		return
	}

	fmt.Println("Trustline, sync_in, and timestamp updated successfully.")

	// Prepare the datagram to send back to the peer
	dg := main.Datagram{
		Command:        main.Server_SetSyncOut,
		XUsername:      ctx.Datagram.YUsername,       // Reverse the usernames for response
		YUsername:      ctx.Datagram.XUsername,
		YServerAddress: main.GetServerAddress(),      // Use the server's address
		Counter:        ctx.Datagram.Counter,         // Copy the existing counter directly
	}

	// Replace explicit signing and sending with the centralized function call
	if err := handlers.SignAndSendDatagram(ctx, &dg); err != nil {
		fmt.Printf("Failed to sign and send datagram: %v\n", err)
		return
	}

	// Add a success message indicating all operations were successful
	fmt.Println("Trustline update and datagram sending completed successfully.")
}
