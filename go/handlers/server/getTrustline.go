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

	// Retrieve the current trustline amount
	trustline, err := main.GetTrustlineOut(ctx.Datagram)
	if err != nil {
		fmt.Printf("Error getting trustline: %v\n", err)
		return
	}

	// Retrieve the current counter value
	counter, err := main.GetCounterOut(ctx.Datagram)
	if err != nil {
		fmt.Printf("Error getting counter: %v\n", err)
		return
	}

	// Prepare a new Datagram for SetTrustline command to be sent to the requesting server
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

	fmt.Println("SetTrustline command sent successfully with current counter.")
}
