package main

import (
    "bytes"
    "crypto/sha256"
    "fmt"
    "os"
    "path/filepath"
)

// verifySignature checks if the signature of the datagram is valid.
func verifySignature(dg Datagram, dir string) error {
    secretKey, err := loadSecretKey(dir)
    if err != nil {
        return fmt.Errorf("error loading secret key: %w", err)
    }

    // Create a byte slice that contains the datagram without the signature
    // and append the secret key to it
    dataWithKey := append(dg[:len(dg)-32], secretKey...) // Assuming the signature is the last 32 bytes

    // Generate the SHA-256 hash
    generatedHash := sha256.Sum256(dataWithKey)

    // Compare the generated hash with the provided signature
    if !bytes.Equal(generatedHash[:], dg.Signature[:]) {
        return fmt.Errorf("signature does not match")
    }

    return nil
}

func loadSecretKey(dir string) ([]byte, error) {
    secretKeyPath := filepath.Join(dir, "secretkey.txt")
    return os.ReadFile(secretKeyPath)
}
