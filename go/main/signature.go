package main

import (
    "bytes"
    "crypto/sha256"
    "fmt"
    "os"
    "path/filepath"
)

// generateSignature computes the SHA-256 signature for the given data.
func generateSignature(data []byte, dir string) ([]byte, error) {
    // Load the secret key from the specified directory.
    secretKeyPath := filepath.Join(dir, "secretkey.txt")
    secretKey, err := os.ReadFile(secretKeyPath)
    if err != nil {
        return nil, fmt.Errorf("error reading secret key: %w", err)
    }

    // Create a byte slice that contains the data without the signature
    dataWithKey := append(data[:len(data)-32], secretKey...)

    // Generate the SHA-256 hash
    generatedHash := sha256.Sum256(dataWithKey)

    return generatedHash[:], nil
}

// VerifySignature checks if the signature of the datagram is valid.
func VerifySignature(dg Datagram, dir string) error {
    data := make([]byte, len(dg)-32) // Adjust size based on the actual struct size
    copy(data, dg[:len(dg)-32])      // Exclude the signature part

    generatedHash, err := generateSignature(data, dir)
    if err != nil {
        return err
    }

    // Compare the generated hash with the provided signature
    if !bytes.Equal(generatedHash, dg.Signature[:]) {
        return fmt.Errorf("signature does not match")
    }

    return nil
}
