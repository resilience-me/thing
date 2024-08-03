package main

import (
    "bytes"
    "crypto/sha256"
    "fmt"
    "os"
    "path/filepath"
)

func verifySignature(dg Datagram, dir string) error {
    secretKey, err := loadSecretKey(dir)
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

func loadSecretKey(dir string) ([]byte, error) {
    secretKeyPath := filepath.Join(dir, "secretkey.txt")
    return os.ReadFile(secretKeyPath)
}
