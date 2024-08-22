package database

import (
    "fmt"
    "os"
    "path/filepath"
)

// loadSecretKey loads the secret key from the specified directory.
func loadSecretKeyFromDir(dir string) ([]byte, error) {
    secretKeyPath := filepath.Join(dir, "secretkey.txt")
    secretKey, err := os.ReadFile(secretKeyPath)
    if err != nil {
        return nil, fmt.Errorf("error reading secret key from %s: %w", secretKeyPath, err)
    }
    return secretKey, nil
}

func LoadSecretKey(username string) ([]byte, error) {
    accountDir := GetAccountDir(username)
    return loadSecretKeyFromDir(accountDir)
}

func LoadPeerSecretKey(username, peerServerAddress, peerUsername string) ([]byte, error) {
    peerDir := GetPeerDir(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
    return loadSecretKeyFromDir(peerDir)
}
