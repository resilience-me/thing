package server

import (
    "encoding/binary"
    "fmt"
    "net"
    "os"
    "path/filepath"
    "strconv"
    "time"

    "resilience/main"
)

// SetTrustline handles setting or updating a trustline from another server's perspective
func SetTrustline(ctx main.HandlerContext) {

    if err := handlers.ValidateServerRequest(ctx); err != nil {
        fmt.Printf("Validation failed: %v\n", err) // Log detailed error
        return
    }

    // Get the peer directory
    peerDir := main.GetPeerDir(ctx.Datagram)

    // Get the trustline directory
    trustlineDir := filepath.Join(peerDir, "trustline")

    // Construct the trustline, counter, and timestamp file paths
    counterInPath := filepath.Join(trustlineDir, "counter_in.txt")
    trustlineInPath := filepath.Join(trustlineDir, "trustline_in.txt")
    timestampPath := filepath.Join(trustlineDir, "sync_timestamp.txt")

    // Load the previous counter value
    prevCounterStr, err := os.ReadFile(counterInPath)
    if err != nil && !os.IsNotExist(err) {
        fmt.Printf("Error reading counter file: %v\n", err)
        return
    }

    prevCounter, err := strconv.ParseUint(string(prevCounterStr), 10, 32) // Parse as uint64 first
    if err != nil {
        fmt.Printf("Error parsing string: %v\n", err)
        return
    }

    // Check the counter
    counter := binary.BigEndian.Uint32(ctx.Datagram.Counter[:])
    if counter <= uint32(prevCounter) {
        fmt.Println("Received counter is not greater than previous counter. Potential replay attack.")
        return
    }

    trustlineAmount := binary.BigEndian.Uint32(ctx.Datagram.Arguments[:4])

    // Write the new trustline amount to the file
    if err := os.WriteFile(trustlineInPath, []byte(fmt.Sprintf("%d", trustlineAmount)), 0644); err != nil {
        fmt.Printf("Error writing trustline to file: %v\n", err)
        return
    }

    // Write the new counter value to the file
    if err := os.WriteFile(counterInPath, []byte(fmt.Sprintf("%d", counter)), 0644); err != nil {
        fmt.Printf("Error writing counter to file: %v\n", err)
        return
    }

    // Write the Unix timestamp to the file
    if err := os.WriteFile(timestampPath, []byte(fmt.Sprintf("%d", time.Now().Unix())), 0644); err != nil {
        fmt.Printf("Error writing timestamp to file: %v\n", err)
        return
    }

    fmt.Println("Trustline, counter, and timestamp updated successfully.")

    // Prepare the datagram to send back to the peer
    dg := main.Datagram{
        Command:       main.Server_SetSyncCounter,
        XUsername:     ctx.Datagram.YUsername,       // Reverse the usernames for response
        YUsername:     ctx.Datagram.XUsername,
        YServerAddress: main.GetServerAddress(),      // Use the server's address
        Counter:       ctx.Datagram.Counter,           // Copy the existing counter directly
    }

    // Sign the datagram
    if err := main.SignDatagram(&dg); err != nil {
        fmt.Printf("Failed to sign datagram: %v\n", err)
        return
    }

    // Send the datagram back to the peer
    _, err = ctx.Conn.WriteToUDP(dg[:], ctx.Addr)
    if err != nil {
        fmt.Printf("Error sending datagram: %v\n", err)
    }
}
