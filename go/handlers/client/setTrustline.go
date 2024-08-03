package handlers

import (
    "encoding/binary"
    "fmt"
    "net"
    "os"
    "path/filepath"

    "resilience/main"
)

// SetTrustline handles setting or updating a trustline from the client's perspective
func SetTrustline(dg main.Datagram, addr *net.UDPAddr) {
    trustlineAmount := binary.BigEndian.Uint32(dg.Arguments[:4])

    accountDir, err := main.GetAccountDir(dg)
    if err != nil {
        fmt.Printf("Error getting account directory: %v\n", err)
        return
    }

    peerDir, err := main.GetPeerDir(dg, accountDir)
    if err != nil {
        fmt.Printf("Error getting peer directory: %v\n", err)
        return
    }

    if err := main.verifySignature(dg, accountDir); err != nil {
        fmt.Printf("Signature verification failed: %v\n", err)
        return
    }

    // Get the trustline directory
    trustlineDir := filepath.Join(peerDir, "trustline")

    // Construct the trustline and counter file paths
    counterOutPath := filepath.Join(trustlineDir, "counter_out.txt")
    trustlineOutPath := filepath.Join(trustlineDir, "trustline_out.txt")

    // Load the previous counter value
    prevCounterStr, err := os.ReadFile(counterOutPath)
    if err != nil && !os.IsNotExist(err) {
        fmt.Printf("Error reading counter file: %v\n", err)
        return
    }
    prevCounter := 0
    if len(prevCounterStr) > 0 {
        prevCounter = int(binary.BigEndian.Uint32(prevCounterStr))
    }

    // Check the counter
    counter := binary.BigEndian.Uint32(dg.Counter[:])
    if int(counter) <= prevCounter {
        fmt.Println("Received counter is not greater than previous counter. Potential replay attack.")
        return
    }

    // Write the new trustline amount to the file
    if err := os.WriteFile(trustlineOutPath, []byte(fmt.Sprintf("%d", trustlineAmount)), 0644); err != nil {
        fmt.Printf("Error writing trustline to file: %v\n", err)
        return
    }

    // Write the new counter value as a string
    counterStr := fmt.Sprintf("%d", counter)
    if err := os.WriteFile(counterOutPath, []byte(counterStr), 0644); err != nil {
        fmt.Printf("Error writing counter to file: %v\n", err)
        return
    }

    fmt.Println("Trustline and counter updated successfully.")
}
