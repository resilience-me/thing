package server_trustlines

import (
	"encoding/binary"
	"fmt"
	"time"
	"resilience/main"
	"resilience/handlers"
)

// SetTimestamp handles updating the sync timestamp for trustlines
func SetTimestamp(ctx main.HandlerContext) {

	if err := handlers.ValidateServerRequest(ctx); err != nil {
		fmt.Printf("Validation failed: %v\n", err) // Log detailed error
		return
	}

	// Retrieve the previous sync_counter_in value using the getter
	prevSyncCounterIn, err := main.GetSyncCounterIn(ctx.Datagram)
	if err != nil {
		fmt.Printf("Error getting previous sync_counter_in: %v\n", err)
		return
	}

	// Get the new counter value from the datagram
	counter := binary.BigEndian.Uint32(ctx.Datagram.Counter[:])

	// Check if the new counter is greater than the previous sync_counter_in
	if counter <= prevSyncCounterIn {
		fmt.Println("Received counter is not greater than previous sync_counter_in. Potential replay attack.")
		return
	}

	// Write the new sync_counter_in value using the setter
	if err := main.SetSyncCounterIn(ctx.Datagram, counter); err != nil {
		fmt.Printf("Error writing sync_counter_in to file: %v\n", err)
		return
	}

	// Retrieve the current timestamp
	timestamp := time.Now().Unix()

	// Write the new timestamp using the setter
	if err := main.SetTimestamp(ctx.Datagram, timestamp); err != nil {
		fmt.Printf("Error writing timestamp to file: %v\n", err)
		return
	}

	fmt.Println("Sync_counter_in and timestamp updated successfully.")
}
