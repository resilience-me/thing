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

func loadClientSecretKey(dg *Datagram) ([]byte, error) {
    keyDir := GetAccountDir(dg)
    return loadSecretKeyFromDir(keyDir)
}

func loadServerSecretKey(dg *Datagram) ([]byte, error) {
    keyDir := GetPeerDir(dg)
    return loadSecretKeyFromDir(keyDir)
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
