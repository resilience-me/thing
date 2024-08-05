package main

import (
    "os"
    "path/filepath"
)

// Initialize datadir only once
var datadir = filepath.Join(os.Getenv("HOME"), "resilience")

// GetAccountDir constructs the account directory path from the datagram
func GetAccountDir(dg *Datagram) (string, error) { // Change to pointer
    username := string(dg.XUsername[:])
    accountDir := filepath.Join(datadir, "accounts", username)

    // Ensure the account directory exists
    if _, err := os.Stat(accountDir); err != nil {
        return "", err
    }

    return accountDir, nil
}

// GetPeerDir constructs the peer directory path from the datagram and checks if it exists.
func GetPeerDir(dg *Datagram, accountDir string) (string, error) { // Change to pointer
    peerUsername := string(dg.YUsername[:])
    peerAddress := string(dg.YServerAddress[:])

    peerDir := filepath.Join(accountDir, "peers", peerAddress, peerUsername)

    // Ensure the peer directory exists
    if _, err := os.Stat(peerDir); err != nil {
        return "", err
    }

    return peerDir, nil
}
