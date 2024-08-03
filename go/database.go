package main

import (
    "fmt"
    "os"
    "path/filepath"
)

// Initialize datadir only once
var datadir = filepath.Join(os.Getenv("HOME"), "ripple")

// GetPeerDir constructs the peer directory path from the datagram and checks if it exists.
func GetPeerDir(dg Datagram) (string, error) {
    username := string(dg.XUsername[:])
    peerUsername := string(dg.YUsername[:])
    peerAddress := string(dg.YServerAddress[:])

    peerDir := filepath.Join(datadir, "accounts", username, "peers", peerAddress, peerUsername)

    // Ensure the peer directory exists
    if _, err := os.Stat(peerDir); err != nil {
        return "", err
    }

    return peerDir, nil
}
