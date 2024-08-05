package server

import (
    "fmt"
    "net"
    "os"
    "path/filepath"
    "encoding/binary"

    "resilience/main"
    "resilience/handlers"
)

// SetSyncCounter handles updating the sync counter from a received context
func SetSyncCounter(ctx main.HandlerContext) {

    if err := handlers.ValidateServerRequest(ctx); err != nil {
        fmt.Printf("Validation failed: %v\n", err) // Log detailed error
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
