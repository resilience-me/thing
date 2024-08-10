package main

import (
    "bytes"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/binary"
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

// verifyHMAC checks the integrity and authenticity of the received buffer
func verifyHMAC(buf []byte, key []byte) bool {
    // The signature is the last 32 bytes of the buffer
    data := buf[:len(buf)-32]
    signature := buf[len(buf)-32:]

    mac := hmac.New(sha256.New, key)
    mac.Write(data)
    expectedMAC := mac.Sum(nil)

    return hmac.Equal(signature, expectedMAC)
}

func loadSecretKey(dg *Datagram) ([]byte, error) {
    var keyDir string
    // Assuming that username, peerAddress, and peerUsername are fields of dg
    if dg.Command & 0x80 == 0 {
        keyDir = filepath.Join(datadir, "accounts", dg.Username)
    } else {
        keyDir = filepath.Join(datadir, "accounts", dg.Username, "peers", dg.PeerServerAddress, dg.PeerUsername)
    }
    return loadSecretKeyFromDir(keyDir)
}

func checkAccountsExist(dg *Datagram) error {
    if err := checkAccountExists(dg); err != nil {
        return fmt.Errorf("error checking account existence: %v", err)
    }
    if err := checkPeerExists(dg); err != nil {
        return fmt.Errorf("error checking peer existence: %v", err)
    }
}

func validateAndParseDatagram(buf []byte) (*Datagram, error) {
    dg := parseDatagram(buf)

    if err := checkAccountsExist(dg); err != nil {
        return nil, err
    }

    secretKey, err := loadSecretKey(dg)
    if err != nil {
        return nil, err
    }
    if !verifyHMAC(buf, secretKey) {
        return nil, errors.New("error verifying HMAC")
    }
    return dg, nil
}
