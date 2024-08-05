package server

import (
    "fmt"
    "net"
    "os"
    "path/filepath"
    "encoding/binary"
    "resilience/main"
)

// SetSyncCounter handles updating the sync counter from a received context
func SetSyncCounter(ctx main.HandlerContext) {
    // Check if the account exists
    if err := main.CheckAccountExists(ctx.Datagram); err != nil {
        fmt.Printf("Error checking account existence: %v\n", err)
        return
    }

    // Check if the peer exists
    if err := main.CheckPeerExists(ctx.Datagram); err != nil {
        fmt.Printf("Error checking peer existence: %v\n", err)
        return
    }

    // Verify the signature
    if err := main.VerifyServerSignature(ctx.Datagram); err != nil {
        fmt.Printf("Signature verification failed: %v\n", err)
        return
    }

    // Get the peer directory using the datagram
    peerDir := main.GetPeerDir(ctx.Datagram)

    // Define the path for sync_counter.txt in the peer directory
    syncCounterPath := filepath.Join(peerDir, "trustline", "sync_counter.txt")

    // Get the current sync counter from the datagram
    syncCounter := binary.BigEndian.Uint32(ctx.Datagram.Counter[:])

    // Write the sync counter to sync_counter.txt
    if err := os.WriteFile(syncCounterPath, []byte(fmt.Sprintf("%d", syncCounter)), 0644); err != nil {
        fmt.Printf("Error writing sync counter to file: %v\n", err)
        return
    }

    fmt.Println("Sync counter updated successfully.")
}
