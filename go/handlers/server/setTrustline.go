package server

import (
    "encoding/binary"
    "fmt"
    "net"
    "os"
    "path/filepath"
    "time"

    "resilience/main"
)

// SetTrustline handles setting or updating a trustline from the server's perspective
func SetTrustline(dg main.Datagram, addr *net.UDPAddr, conn *net.UDPConn) {
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

    if err := main.VerifySignature(dg, peerDir); err != nil {
        fmt.Printf("Signature verification failed: %v\n", err)
        return
    }

    // Get the trustline directory
    trustlineDir := filepath.Join(peerDir, "trustline")

    // Construct the trustline, counter and timestamp file paths
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
    counter := binary.BigEndian.Uint32(dg.Counter[:])
    if counter <= uint32(prevCounter) {
        fmt.Println("Received counter is not greater than previous counter. Potential replay attack.")
        return
    }

    // Write the new trustline amount to the file
    if err := os.WriteFile(trustlineInPath, []byte(fmt.Sprintf("%d", trustlineAmount)), 0644); err != nil {
        fmt.Printf("Error writing trustline to file: %v\n", err)
        return
    }

    // Write the new counter value to the file directly
    if err := os.WriteFile(counterInPath, []byte(fmt.Sprintf("%d", counter)), 0644); err != nil {
        fmt.Printf("Error writing counter to file: %v\n", err)
        return
    }

    // Write the Unix timestamp directly to the file at the end if everything is successful
    if err := os.WriteFile(timestampPath, []byte(fmt.Sprintf("%d", time.Now().Unix())), 0644); err != nil {
        fmt.Printf("Error writing timestamp to file: %v\n", err)
        return
    }

    fmt.Println("Trustline, counter and timestamp updated successfully.")

    // Prepare the response datagram
    responseDg := main.Datagram{
        Command:       main.Server_SetSyncCounter,
        XUsername:     dg.YUsername,        // Reverse the usernames for response
        YUsername:     dg.XUsername,
        YServerAddress: main.GetServerAddress(), // Use the server's address
        Counter:       dg.Counter,           // Copy the existing counter directly
    }

    // Sign the response datagram
    if err := main.SignDatagram(&responseDg, peerDir); err != nil {
        fmt.Printf("Error signing response datagram: %v\n", err)
        return
    }

    // Send the response datagram back to the peer
    _, err = conn.WriteToUDP(responseDg[:], addr)
    if err != nil {
        fmt.Printf("Error sending response datagram: %v\n", err)
    }
}
