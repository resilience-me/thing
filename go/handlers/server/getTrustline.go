package server

import (
	"encoding/binary"
	"fmt"
	"resilience/main"
	"resilience/handlers"
)

// GetTrustline handles the request to get the current trustline amount from another server
func GetTrustline(ctx main.HandlerContext) {
	// Validate the server request
	if err := handlers.ValidateServerRequest(ctx); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
		return
	}

	// Retrieve the current counter value
	counter, err := main.GetCounter(ctx.Datagram)
	if err != nil {
		fmt.Printf("Error getting counter: %v\n", err)
		return
	}

	// Retrieve the current sync_out value
	syncOut, err := main.GetSyncOut(ctx.Datagram)
	if err != nil {
		fmt.Printf("Error getting sync_out: %v\n", err)
		return
	}

	// Check if the server is synced
	if counter == syncOut {
		// Synced, send SetTrustlineSyncTimestamp command
		syncCounterOut, err := main.GetSyncCounterOut(ctx.Datagram)
		if err != nil {
			fmt.Printf("Error getting sync_counter_out: %v\n", err)
			return
		}

		dg := main.Datagram{
			Command: main.Server_SetTrustlineSyncTimestamp,
			XUsername: ctx.Datagram.YUsername, // Assume reverse usernames for response
			YUsername: ctx.Datagram.XUsername,
			YServerAddress: main.GetServerAddress(), // Server's own address
		}

		// Set the sync_counter_out as the counter in the Datagram
		binary.BigEndian.PutUint32(dg.Counter[:], syncCounterOut)

		// Use the handlers.SignAndSendDatagram to sign and send the datagram
		if err := handlers.SignAndSendDatagram(ctx, &dg); err != nil {
			fmt.Printf("Failed to sign and send datagram: %v\n", err)
			return
		}

		fmt.Println("SetTrustlineSyncTimestamp command sent successfully.")
	} else {
		// Not synced, send SetTrustline command
		// Retrieve the current trustline amount
		trustline, err := main.GetTrustlineOut(ctx.Datagram)
		if err != nil {
			fmt.Printf("Error getting trustline: %v\n", err)
			return
		}

		dg := main.Datagram{
			Command: main.Server_SetTrustline,
			XUsername: ctx.Datagram.YUsername, // Assume reverse usernames for response
			YUsername: ctx.Datagram.XUsername,
			YServerAddress: main.GetServerAddress(), // Server's own address
		}

		// Set the trustline amount in the arguments section of the Datagram
		binary.BigEndian.PutUint32(dg.Arguments[:4], trustline)
		// Set the current counter in the Datagram
		binary.BigEndian.PutUint32(dg.Counter[:], counter)

		// Use the handlers.SignAndSendDatagram to sign and send the datagram
		if err := handlers.SignAndSendDatagram(ctx, &dg); err != nil {
			fmt.Printf("Failed to sign and send datagram: %v\n", err)
			return
		}

		fmt.Println("SetTrustline command sent successfully.")
	}
}
