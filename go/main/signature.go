package main

import (
    "bytes"
    "crypto/sha256"
    "fmt"
    "os"
    "path/filepath"
)

// loadSecretKey loads the secret key from the specified directory.
func loadSecretKey(dir string) ([]byte, error) {
    secretKeyPath := filepath.Join(dir, "secretkey.txt")
    secretKey, err := os.ReadFile(secretKeyPath)
    if err != nil {
        return nil, fmt.Errorf("error reading secret key from %s: %w", secretKeyPath, err)
    }
    return secretKey, nil
}

// generateSignature computes the SHA-256 signature for the given data and secret key.
func generateSignature(data []byte, secretKey []byte) ([]byte, error) {    
    // Create a byte slice that contains the data without the signature
    dataWithKey := append(data[:len(data)-32], secretKey...)

    // Generate the SHA-256 hash
    generatedHash := sha256.Sum256(dataWithKey)

    return generatedHash[:], nil
}

// SignDatagram signs the given Datagram by generating a signature.
func SignDatagram(dg *Datagram) error {
    // Get the peer directory based on the datagram
    peerDir := GetPeerDir(dg)

    // Load the secret key
    secretKey, err := loadSecretKey(peerDir)
    if err != nil {
        return fmt.Errorf("failed to load secret key in SignDatagram: %w", err)
    }

    // Call generateSignature directly with the Datagram's byte representation and the secret key
    signature, err := generateSignature((*dg)[:], secretKey)
    if err != nil {
        return fmt.Errorf("failed to generate signature for Datagram: %w", err)
    }

    // Copy the generated signature into the datagram's signature field
    copy(dg.Signature[:], signature)

    return nil
}

// SignResponseDatagram signs the given ResponseDatagram by generating a signature.
func SignResponseDatagram(rd *ResponseDatagram, username string) error {
    // Construct the account directory path from the username
    accountDir := filepath.Join(datadir, "accounts", username)

    // Load the secret key
    secretKey, err := loadSecretKey(accountDir)
    if err != nil {
        return fmt.Errorf("failed to load secret key in SignResponseDatagram: %w", err)
    }

    // Call generateSignature directly with the ResponseDatagram's byte representation and the secret key
    signature, err := generateSignature((*rd)[:], secretKey)
    if err != nil {
        return fmt.Errorf("failed to generate signature for ResponseDatagram: %w", err)
    }

    // Copy the generated signature into the response datagram's signature field
    copy(rd.Signature[:], signature)

    return nil
}

// VerifySignature checks if the signature of the datagram is valid.
func verifySignature(dg *Datagram, dir string) error {
    // Load the secret key
    secretKey, err := loadSecretKey(dir)
    if err != nil {
        return fmt.Errorf("failed to load secret key for verification: %w", err)
    }

    // Generate the expected signature based on the entire datagram
    generatedHash, err := generateSignature(dg[:], secretKey)
    if err != nil {
        return fmt.Errorf("failed to generate signature for verification: %w", err)
    }

    // Compare the generated hash with the provided signature
    if !bytes.Equal(generatedHash, dg.Signature[:]) {
        return fmt.Errorf("signature does not match")
    }

    return nil
}

// VerifyClientSignature verifies the client's signature of the datagram.
func VerifyClientSignature(dg *Datagram) error {
    accountDir := GetAccountDir(dg)
    return verifySignature(dg, accountDir)
}

// VerifyServerSignature verifies the server's signature of the datagram.
func VerifyServerSignature(dg *Datagram) error {
    peerDir := GetPeerDir(dg)
    return verifySignature(dg, peerDir)
}
