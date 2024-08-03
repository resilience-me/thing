package main

import (
    "bytes"
    "crypto/sha256"
    "fmt"
    "os"
    "path/filepath"
)

// generateSignature computes the SHA-256 signature for the given datagram.
func generateSignature(dg Datagram, dir string) ([]byte, error) {
    secretKey, err := loadSecretKey(dir)
    if err != nil {
        return nil, fmt.Errorf("error loading secret key: %w", err)
    }

    // Create a byte slice that contains the datagram without the signature
    dataWithKey := append(dg[:len(dg)-32], secretKey...) // Assuming the signature is the last 32 bytes

    // Generate the SHA-256 hash
    generatedHash := sha256.Sum256(dataWithKey)

    return generatedHash[:], nil
}

// signDatagram signs the given datagram by generating a signature.
func signDatagram(dg *Datagram, dir string) error {
    signature, err := generateSignature(*dg, dir)
    if err != nil {
        return err
    }

    // Copy the generated signature into the datagram's signature field
    copy(dg.Signature[:], signature)

    return nil
}

// verifySignature checks if the signature of the datagram is valid.
func verifySignature(dg Datagram, dir string) error {
    // Generate the expected signature based on the datagram
    expectedSignature, err := generateSignature(dg, dir)
    if err != nil {
        return err
    }

    // Compare the generated hash with the provided signature
    if !bytes.Equal(expectedSignature, dg.Signature[:]) {
        return fmt.Errorf("signature does not match")
    }

    return nil
}


func loadSecretKey(dir string) ([]byte, error) {
    secretKeyPath := filepath.Join(dir, "secretkey.txt")
    return os.ReadFile(secretKeyPath)
}
