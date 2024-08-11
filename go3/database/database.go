package database

import (
    "os"
    "path/filepath"
)

var datadir = filepath.Join(os.Getenv("HOME"), "ripple")

// GetAccountDir constructs the account directory path from the datagram
func GetAccountDir(dg *Datagram) string {
    return filepath.Join(datadir, "accounts", dg.Username)
}

// GetPeerDir constructs the peer directory path from the datagram and returns it
func GetPeerDir(dg *Datagram) string {
    accountDir := GetAccountDir(dg)
    return filepath.Join(accountDir, "peers", dg.PeerServerAddress, dg.PeerUsername)
}

// GetTrustlineDir constructs the trustline directory path from the datagram and returns it.
func GetTrustlineDir(dg *Datagram) string {
    peerDir := GetPeerDir(dg)
    return filepath.Join(peerDir, "trustline")
}

// checkDirExists checks if a specific directory exists.
func checkDirExists(dirPath string) (bool, error) {
    // Use os.Stat to attempt to retrieve the directory information
    if _, err := os.Stat(dirPath); err != nil {
        if os.IsNotExist(err) {
            // The directory does not exist
            return false, nil
        }
        // Return false along with the error encountered during Stat
        return false, err
    }
    // The directory exists
    return true, nil
}

// CheckAccountExists checks if the account directory exists
func CheckAccountExists(dg *Datagram) (bool, error) {
    accountDir := GetAccountDir(dg)
    // Ensure the account directory exists
    return checkDirExists(accountDir)
}

// CheckPeerExists checks if the peer directory exists
func CheckPeerExists(dg *Datagram) (bool, error) {
    peerDir := GetPeerDir(dg)
    // Ensure the peer directory exists
    return checkDirExists(peerDir)
}

// checkUserAndPeerExist checks if both account and peer directories exist for the given datagram.
func checkUserAndPeerExist(dg *Datagram) (bool, error) {
    // Check if account directory exists
    if exists, err := CheckAccountExists(dg); err != nil {
        return false, fmt.Errorf("error checking account existence for user '%s': %v", dg.Username, err)
    } else if !exists {
        return false, fmt.Errorf("account directory does not exist for user '%s'", dg.Username)
    }

    // Check if peer directory exists
    if exists, err = CheckPeerExists(dg); err != nil {
        return false, fmt.Errorf("error checking peer existence for server '%s' and user '%s': %v", dg.PeerServerAddress, dg.PeerUsername, err)
    } else if !exists {
        return false, fmt.Errorf("peer directory does not exist for server '%s' and user '%s'", dg.PeerServerAddress, dg.PeerUsername)
    }

    return true, nil
}
