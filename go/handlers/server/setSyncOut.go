package server

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

	// Get the current sync_out value from the datagram
	syncOut := binary.BigEndian.Uint32(ctx.Datagram.Counter[:])

	// Write the sync_out value using the setter
	if err := main.SetSyncOut(ctx.Datagram, syncOut); err != nil {
		fmt.Printf("Error writing sync_out to file: %v\n", err)
		return
	}

	fmt.Println("Sync_out updated successfully.")
}
