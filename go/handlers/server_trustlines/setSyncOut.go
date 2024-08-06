package server_trustlines

import (
	"encoding/binary"
	"fmt"
	"resilience/main"
	"resilience/handlers"
)

// SetSyncOut handles updating the sync_out counter from a received context
func SetSyncOut(ctx main.HandlerContext) {

	if err := handlers.ValidateServerRequest(ctx); err != nil {
		fmt.Printf("Validation failed: %v\n", err) // Log detailed error
		return
	}

	// Retrieve the previous sync_out value using the getter
	prevSyncOut, err := main.GetSyncOut(ctx.Datagram)
	if err != nil {
		fmt.Printf("Error getting previous sync_out: %v\n", err)
		return
	}

	// Get the new sync_out value from the datagram
	syncOut := binary.BigEndian.Uint32(ctx.Datagram.Counter[:])

	// Check if the new sync_out is greater than the previous sync_out
	if syncOut <= prevSyncOut {
		fmt.Println("Received sync_out is not greater than previous sync_out. Potential replay attack.")
		return
	}

	// Write the new sync_out value using the setter
	if err := main.SetSyncOut(ctx.Datagram, syncOut); err != nil {
		fmt.Printf("Error writing sync_out to file: %v\n", err)
		return
	}

	fmt.Println("Sync_out updated successfully.")
}
