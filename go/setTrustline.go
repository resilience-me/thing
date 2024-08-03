package main

import (
    "encoding/binary"
    "fmt"
    "net"
    "os"
    "path/filepath"
)

func setTrustline(dg Datagram, addr *net.UDPAddr) {
    trustlineAmount := binary.BigEndian.Uint32(dg.Arguments[:4])

    peerDir, err := GetPeerDir(dg)
    if err != nil {
        fmt.Printf("Error getting peer directory: %v\n", err)
        return
    }

    if err := verifySignature(dg, peerDir); err != nil {
        fmt.Printf("Signature verification failed: %v\n", err)
        return
    }

    // Construct the trustline and counter file paths
    trustlineDir := filepath.Join(peerDir, "trustline")
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
    incomingCounter := binary.BigEndian.Uint32(dg.Counter[:])
    if int(incomingCounter) <= prevCounter {
        fmt.Println("Received counter is not greater than previous counter. Potential replay attack.")
        return
    }

    // Write the new trustline amount to the file
    if err := os.WriteFile(trustlineOutPath, []byte(fmt.Sprintf("%d", trustlineAmount)), 0644); err != nil {
        fmt.Printf("Error writing trustline to file: %v\n", err)
        return
    }

    // Increment and write the new counter value
    newCounter := incomingCounter + 1
    counterData := make([]byte, 4)
    binary.BigEndian.PutUint32(counterData, newCounter)
    if err := os.WriteFile(counterOutPath, counterData, 0644); err != nil {
        fmt.Printf("Error writing counter to file: %v\n", err)
        return
    }

    fmt.Println("Trustline and counter updated successfully.")
}
