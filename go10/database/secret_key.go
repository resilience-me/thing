package database

import "fmt"

// loadSecretKeyFromDir loads the secret key from the specified directory.
func loadSecretKeyFromDir(dir string) ([]byte, error) {
    secretKey, err := ReadFile(dir, "secretkey.txt")
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
