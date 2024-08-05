package main

import (
    "os"
    "path/filepath"
)

// Initialize datadir only once
var datadir = filepath.Join(os.Getenv("HOME"), "resilience")

// GetAccountDir constructs the account directory path from the datagram
func GetAccountDir(dg *Datagram) string { // Change to pointer
    username := string(dg.XUsername[:])
    accountDir := filepath.Join(datadir, "accounts", username)

    return accountDir
}

// GetPeerDir constructs the peer directory path from the datagram and returns it
func GetPeerDir(dg *Datagram, accountDir string) string { // Change to pointer
    peerUsername := string(dg.YUsername[:])
    peerAddress := string(dg.YServerAddress[:])

    peerDir := filepath.Join(accountDir, "peers", peerAddress, peerUsername)

    return peerDir
}

// CheckAccountExists checks if the account directory exists
func CheckAccountExists(accountDir string) error {
    // Ensure the account directory exists
    if err := os.Stat(accountDir); err != nil {
        return err
    }
    return nil
}

// CheckPeerExists checks if the peer directory exists
func CheckPeerExists(peerDir string) error {
    // Ensure the peer directory exists
    if err := os.Stat(peerDir); err != nil {
        return err
    }
    return nil
}
