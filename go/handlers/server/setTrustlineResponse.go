package server

import (
    "fmt"
    "net"
    "os"
    "path/filepath"
    "encoding/binary"
    "strconv"
    "resilience/main"
)

func SetTrustlineResponse(dg main.Datagram, addr *net.UDPAddr) {
    // Get the account directory
    accountDir, err := main.GetAccountDir(dg)
    if err != nil {
        fmt.Printf("Error getting account directory: %v\n", err)
        return
    }

    // Get the peer directory
    peerDir, err := main.GetPeerDir(dg, accountDir)
    if err != nil {
        fmt.Printf("Error getting peer directory: %v\n", err)
        return
    }

    // Verify the signature using peerDir
    if err := main.VerifySignature(dg, peerDir); err != nil {
        fmt.Printf("Signature verification failed: %v\n", err)
        return
    }

    // Define the path for sync_counter.txt in peerDir
    syncCounterPath := filepath.Join(peerDir, "trustline", "sync_counter.txt")

    // Get the current sync counter from the datagram
    syncCounter := binary.BigEndian.Uint32(dg.Counter[:])

    // Write the sync counter to sync_counter.txt
    if err := os.WriteFile(syncCounterPath, []byte(fmt.Sprintf("%d", syncCounter)), 0644); err != nil {
        fmt.Printf("Error writing sync counter to file: %v\n", err)
        return
    }

    fmt.Println("Sync counter updated successfully.")
}
