package server

import (
	"encoding/binary"
	"fmt"
	"time"
	"resilience/main"
	"resilience/handlers"
)

// SetTrustlineSyncTimestamp handles updating the sync timestamp for trustlines
func SetTrustlineSyncTimestamp(ctx main.HandlerContext) {

	if err := handlers.ValidateServerRequest(ctx); err != nil {
		fmt.Printf("Validation failed: %v\n", err) // Log detailed error
		return
	}

	// Retrieve the current timestamp
	timestamp := time.Now().Unix()

	// Write the new timestamp using the setter
	if err := main.SetTimestamp(ctx.Datagram, timestamp); err != nil {
		fmt.Printf("Error writing timestamp to file: %v\n", err)
		return
	}

	fmt.Println("Timestamp updated successfully.")
}
