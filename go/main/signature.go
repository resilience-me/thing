package main

import (
    "bytes"
    "crypto/sha256"
    "fmt"
    "os"
    "path/filepath"
)

// loadSecretKey loads the secret key from the specified directory.
func loadSecretKey(dir string) ([32]byte, error) {
    secretKeyPath := filepath.Join(dir, "secretkey.txt")
    secretKeyBytes, err := os.ReadFile(secretKeyPath)
    if err != nil {
        return [32]byte{}, fmt.Errorf("error reading secret key from %s: %w", secretKeyPath, err)
    }
    // Convert secretKeyBytes to [32]byte
    var secretKey [32]byte
    copy(secretKey[:], secretKeyBytes) // Copy the contents to the fixed-size array

    return secretKey, nil
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

    // Write the secret key into the signature field
    dg.Signature = secretKey

    // Generate the signature using the current datagram (with the secret key in the signature field)
    generatedHash := sha256.Sum256(*dg)

    // Replace the signature field with the generated hash
    dg.Signature = generatedHash

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

    rd.Signature = secretKey

    // Generate the signature using the current datagram (with the secret key in the signature field)
    signature := sha256.Sum256(*rd)

    // Copy the generated signature into the response datagram's signature field
    rd.Signature = signature

    return nil
}

// VerifySignature checks if the signature of the datagram is valid.
func verifySignature(dg *Datagram, dir string) error {
    // Load the secret key
    secretKey, err := loadSecretKey(dir)
    if err != nil {
        return fmt.Errorf("failed to load secret key for verification: %w", err)
    }
    buf := *dg;
    buf.Signature = secretKey
    // Generate the expected signature based on the entire datagram
    generatedHash := sha256.Sum256(buf)

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
