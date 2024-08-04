package main

import (
    "bytes"
    "crypto/sha256"
    "fmt"
    "os"
    "path/filepath"
)

// generateSignature computes the SHA-256 signature for the given data.
func GenerateSignature(data []byte, dir string) ([]byte, error) {
    // Load the secret key from the specified directory.
    secretKeyPath := filepath.Join(dir, "secretkey.txt")
    secretKey, err := os.ReadFile(secretKeyPath)
    if err != nil {
        return nil, fmt.Errorf("error reading secret key: %w", err)
    }

    // Check if data is at least 32 bytes long
    if len(data) < 32 {
        return nil, fmt.Errorf("data must be at least 32 bytes long")
    }
    
    // Create a byte slice that contains the data without the signature
    dataWithKey := append(data[:len(data)-32], secretKey...)

    // Generate the SHA-256 hash
    generatedHash := sha256.Sum256(dataWithKey)

    return generatedHash[:], nil
}

// VerifySignature checks if the signature of the datagram is valid.
func VerifySignature(dg Datagram, dir string) error {

    // Generate the expected signature based on the entire datagram
    generatedHash, err := GenerateSignature(dg[:], dir)
    if err != nil {
        return err
    }

    // Compare the generated hash with the provided signature
    if !bytes.Equal(generatedHash, dg.Signature[:]) {
        return fmt.Errorf("signature does not match")
    }

    return nil
}
