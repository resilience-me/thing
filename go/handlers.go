package main

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "net"
    "os"
    "path/filepath"
    "strconv"
)

func setTrustline(dg Datagram, addr *net.UDPAddr) {
    trustlineAmount := binary.BigEndian.Uint32(dg.Arguments[:4])
    username := string(dg.XUsername[:]) // Assuming XUsername is null-terminated ASCII
    peerUsername := string(dg.YUsername[:])
    peerAddress := string(dg.YServerAddress[:])

    // Construct the file paths
    datadir := filepath.Join(os.Getenv("HOME"), "ripple")
    peerDir := filepath.Join(datadir, "accounts", username, "peers", peerAddress, peerUsername)

    // Check if the peer directory exists
    if _, err := os.Stat(peerDir); os.IsNotExist(err) {
        fmt.Println("Peer directory does not exist:", peerDir)
        return
    }

    // Reading the current counter from file
    counterPath := filepath.Join(peerDir, "counter_out.txt")
    file, err := os.Open(counterPath)
    if err != nil {
        fmt.Println("Failed to open counter file:", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Scan() // Read the first line
    currentCounter, _ := strconv.Atoi(scanner.Text())
    fileCounter := binary.BigEndian.Uint32(dg.Counter[:])
    if fileCounter <= uint32(currentCounter) {
        fmt.Println("Counter is not valid or replay attack detected")
        return
    }

    // Writing the new trustline to file
    trustlinePath := filepath.Join(peerDir, "trustline_out.txt")
    err = os.WriteFile(trustlinePath, []byte(strconv.Itoa(int(trustlineAmount))), 0644)
    if err != nil {
        fmt.Println("Failed to write trustline to file:", err)
        return
    }

    // Update the counter
    err = os.WriteFile(counterPath, []byte(strconv.Itoa(int(fileCounter))), 0644)
    if err != nil {
        fmt.Println("Failed to update counter file:", err)
        return
    }

    fmt.Println("Trustline updated successfully for", peerUsername, "to", trustlineAmount)
}
