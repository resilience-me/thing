package server

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"resilience/main"
	"resilience/handlers"
	"strconv"
)

// GetTrustline handles the request to get the current trustline amount from another server
func GetTrustline(ctx main.HandlerContext) {
	// Validate the server request
	if err := handlers.ValidateServerRequest(ctx); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
		return
	}

	// Get the trustline directory
	trustlineDir := main.GetTrustlineDir(ctx.Datagram)
	trustlinePath := filepath.Join(trustlineDir, "trustline_out.txt")
	counterOutPath := filepath.Join(trustlineDir, "counter_out.txt")

	// Read the current trustline amount
	trustlineAmountBytes, err := os.ReadFile(trustlinePath)
	if err != nil {
		fmt.Printf("Error reading trustline amount: %v\n", err)
		return
	}

	trustlineAmount, err := strconv.ParseUint(string(trustlineAmountBytes), 10, 32)
	if err != nil {
		fmt.Printf("Error parsing trustline amount: %v\n", err)
		return
	}

	// Read the current counter value (alphanumeric)
	counterStr, err := os.ReadFile(counterOutPath)
	if err != nil {
		fmt.Printf("Error reading counter: %v\n", err)
		return
	}

	counter, err := strconv.ParseUint(string(counterStr), 10, 32)
	if err != nil {
		fmt.Printf("Error parsing counter: %v\n", err)
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
	binary.BigEndian.PutUint32(dg.Arguments[:4], uint32(trustlineAmount))
	// Set the current counter in the Datagram
	binary.BigEndian.PutUint32(dg.Counter[:], uint32(counter))

	// Use the handlers.SignAndSendDatagram to sign and send the datagram
	if err := handlers.SignAndSendDatagram(ctx, &dg); err != nil {
		fmt.Printf("Failed to sign and send datagram: %v\n", err)
		return
	}

	fmt.Println("SetTrustline command sent successfully with current counter.")
}
