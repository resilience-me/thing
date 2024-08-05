package server

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"path/filepath"
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

	// Construct the path for the trustline information
	peerDir := main.GetPeerDir(ctx.Datagram)
	trustlinePath := filepath.Join(peerDir, "trustline", "trustline_out.txt")

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

	// Prepare a new Datagram for SetTrustline command to be sent to the requesting server
	dg := main.Datagram{
		Command: main.Server_SetTrustline,
		XUsername: ctx.Datagram.YUsername, // Assume reverse usernames for response
		YUsername: ctx.Datagram.XUsername,
		YServerAddress: main.GetServerAddress(), // Server's own address
	}

	// Set the trustline amount in the arguments section of the Datagram
	binary.BigEndian.PutUint32(dg.Arguments[:4], uint32(trustlineAmount))

	// Sign the datagram to ensure integrity and authenticity
	if err := main.SignDatagram(&dg); err != nil {
		fmt.Printf("Failed to sign datagram: %v\n", err)
		return
	}

	// Send the datagram
	_, err = ctx.Conn.WriteToUDP(dg[:], ctx.Addr)
	if err != nil {
		fmt.Printf("Error sending SetTrustline command: %v\n", err)
	}
}
