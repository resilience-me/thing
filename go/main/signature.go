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
    secretKey err := os.ReadFile(secretKeyPath)
    if err != nil {
        return nil, fmt.Errorf("error reading secret key from %s: %w", secretKeyPath, err)
    }

    return secretKey, nil
}

func generateDatagramSignature(dg *Datagram, secretKey []byte) [32]byte {
    var dataWithKey []byte
    dataWithKey = append(dataWithKey, dg.Command)
    dataWithKey = append(dataWithKey, dg.XUsername[:]...)
    dataWithKey = append(dataWithKey, dg.YUsername[:]...)
    dataWithKey = append(dataWithKey, dg.YServerAddress[:]...)
    dataWithKey = append(dataWithKey, dg.Arguments[:]...)
    dataWithKey = append(dataWithKey, dg.Counter[:]...)
    dataWithKey = append(dataWithKey, secretKey...)
    
    generatedHash := sha256.Sum256(dataWithKey)
    return generatedHash
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

    // Replace the signature field with the generated hash
    dg.Signature = generateDatagramSignature(dg, secretKey)

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

    var dataWithKey []byte
    dataWithKey = append(dataWithKey, rd.Nonce[:]...)
    dataWithKey = append(dataWithKey, rd.Result[:]...)
    dataWithKey = append(dataWithKey, secretKey...)
    
    signature := sha256.Sum256(dataWithKey)

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
    // Generate the expected signature based on the datagram
    generatedHash := generateDatagramSignature(dg, secretKey)

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
