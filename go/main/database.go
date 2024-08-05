package main

import (
    "os"
    "path/filepath"
)

// Initialize datadir only once
var datadir = filepath.Join(os.Getenv("HOME"), "resilience")

// GetAccountDir constructs the account directory path from the datagram
func GetAccountDir(dg *Datagram) string {
    username := string(dg.XUsername[:])

    return filepath.Join(datadir, "accounts", username)
}

// GetPeerDir constructs the peer directory path from the datagram and returns it
func GetPeerDir(dg *Datagram) string {
    username := string(dg.XUsername[:])    
    peerUsername := string(dg.YUsername[:])
    peerAddress := string(dg.YServerAddress[:])

    return filepath.Join(datadir, "accounts", username, "peers", peerAddress, peerUsername)
}

// CheckAccountExists checks if the account directory exists
func CheckAccountExists(dg *Datagram) error {
    accountDir := GetAccountDir(dg)
    // Ensure the account directory exists
    if err := os.Stat(accountDir); err != nil {
        return err
    }
    return nil
}

// CheckPeerExists checks if the peer directory exists
func CheckPeerExists(dg *Datagram) error {
    peerDir := GetPeerDir(dg)
    // Ensure the peer directory exists
    if err := os.Stat(peerDir); err != nil {
        return err
    }
    return nil
}
