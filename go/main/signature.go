package main

import (
    "bytes"
    "crypto/sha256"
    "fmt"
    "os"
    "path/filepath"
)

func verifySignature(dg Datagram, peerDir string) error {
    secretKey, err := loadSecretKey(peerDir)
    if err != nil {
        return fmt.Errorf("error loading secret key: %w", err)
    }

    dataWithKey := append(dg.Command[:len(dg.Command)-len(dg.Signature)], secretKey...)
    generatedHash := sha256.Sum256(dataWithKey)

    if !bytes.Equal(generatedHash[:], dg.Signature[:]) {
        return fmt.Errorf("signature does not match")
    }

    return nil
}

func loadSecretKey(peerDir string) ([]byte, error) {
    secretKeyPath := filepath.Join(peerDir, "secretkey.txt")
    return os.ReadFile(secretKeyPath)
}
