package database

import "fmt"

// loadSecretKeyFromDir loads the secret key from the specified directory.
func loadSecretKeyFromDir(dir string) ([]byte, error) {
    secretKey, err := ReadFile(dir, "secretkey.txt")
    if err != nil {
        return nil, fmt.Errorf("error reading secret key: %w", err)
    }
    return secretKey, nil
}

// LoadSecretKey loads the secret key for the given username.
func LoadSecretKey(username string) ([]byte, error) {
    accountDir := GetAccountDir(username)
    return loadSecretKeyFromDir(accountDir)
}

// LoadPeerSecretKey loads the peer's secret key.
func LoadPeerSecretKey(username, peerServerAddress, peerUsername string) ([]byte, error) {
    peerDir := GetPeerDir(username, peerServerAddress, peerUsername)
    return loadSecretKeyFromDir(peerDir)
}
