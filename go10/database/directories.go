package database

import (
    "os"
    "path/filepath"
    "ripple/types"
)

// GetAccountDir constructs the account directory path from a username and returns it
func GetAccountDir(username string) string {
    datadir := main.GetDataDir()
    return filepath.Join(datadir, "accounts", username)
}

// GetPeerDir constructs the peer directory path from a username, peer server address and peer username and returns it
func GetPeerDir(username, peerServerAddress, peerUsername string) string {
    accountDir := GetAccountDir(username)
    return filepath.Join(accountDir, "peers", peerServerAddress, peerUsername)
}

// GetTrustlineDir constructs the trustline directory path from  a username, peer server address and peer username and returns it
func GetTrustlineDir(username, peerServerAddress, peerUsername string) string {
    peerDir := GetPeerDir(username, peerServerAddress, peerUsername)
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

// CheckPeerExists checks if the peer directory exists
func CheckPeerExists(dg *types.Datagram) (bool, error) {
    peerDir := GetPeerDir(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
    // Ensure the peer directory exists
    return checkDirExists(peerDir)
}
