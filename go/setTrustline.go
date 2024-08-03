package main

import (
    "encoding/binary"
    "fmt"
    "net"
    "os"
    "path/filepath"
    "yourapp/data" // Adjust the import path based on your actual module setup
)

func setTrustline(dg data.Datagram, addr *net.UDPAddr) {
    trustlineAmount := binary.BigEndian.Uint32(dg.Arguments[:4])

    peerDir, err := data.GetPeerDir(dg)
    if err != nil {
        fmt.Printf("Error getting peer directory: %v\n", err)
        return
    }

    // Load the secret key
    secretKeyPath := filepath.Join(peerDir, "secretkey.txt")
    secretKey, err := os.ReadFile(secretKeyPath)
    if err != nil {
        fmt.Printf("Error reading secret key: %v\n", err)
        return
    }

    // Signature verification would go here (omitted for brevity)

    // Construct the trustline and counter file paths
    counterOutPath := filepath.Join(peerDir, "counter_out.txt")
    trustlineOutPath := filepath.Join(peerDir, "trustline_out.txt")

    // Load the previous counter value
    prevCounterStr, err := os.ReadFile(counterOutPath);
    
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
