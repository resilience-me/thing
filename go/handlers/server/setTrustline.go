package server

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

    // Get the trustline directory
    trustlineDir := main.GetTrustlineDir(ctx.Datagram)

    // Retrieve the previous counter value using the new getter
    prevCounter, err := main.GetCounterIn(ctx.Datagram)
    if err != nil {
        fmt.Printf("Error getting previous counter: %v\n", err)
        return
    }

    // Check the counter
    counter := binary.BigEndian.Uint32(ctx.Datagram.Counter[:])
    if counter <= prevCounter {
        fmt.Println("Received counter is not greater than previous counter. Potential replay attack.")
        return
    }

    // Retrieve the trustline amount from the Datagram
    trustlineAmount := binary.BigEndian.Uint32(ctx.Datagram.Arguments[:4])

    // Write the new trustline amount to the file
    trustlineInPath := filepath.Join(trustlineDir, "trustline_in.txt")
    if err := os.WriteFile(trustlineInPath, []byte(fmt.Sprintf("%d", trustlineAmount)), 0644); err != nil {
        fmt.Printf("Error writing trustline to file: %v\n", err)
        return
    }

    // Write the new counter value to the file
    counterInPath := filepath.Join(trustlineDir, "counter_in.txt")
    if err := os.WriteFile(counterInPath, []byte(fmt.Sprintf("%d", counter)), 0644); err != nil {
        fmt.Printf("Error writing counter to file: %v\n", err)
        return
    }

    // Write the Unix timestamp to the file
    timestampPath := filepath.Join(trustlineDir, "sync_timestamp.txt")
    if err := os.WriteFile(timestampPath, []byte(fmt.Sprintf("%d", time.Now().Unix())), 0644); err != nil {
        fmt.Printf("Error writing timestamp to file: %v\n", err)
        return
    }

    fmt.Println("Trustline, counter, and timestamp updated successfully.")

    // Prepare the datagram to send back to the peer
    dg := main.Datagram{
        Command:        main.Server_SetSyncCounter,
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
