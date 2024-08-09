package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/hmac"
    "crypto/sha256"
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "strings"
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

// authenticateDatagram authenticates the datagram using HMAC.
func authenticateDatagram(datagram []byte, key []byte) ([]byte, error) {
    if len(datagram) < 422 { // Ensure the datagram meets the minimum length requirement.
        return nil, errors.New("datagram too short")
    }

    // Extract the HMAC from the end of the datagram and separate the data part.
    data, hmacSent := datagram[:len(datagram)-32], datagram[len(datagram)-32:]

    // Inline HMAC verification logic
    mac := hmac.New(sha256.New, key)
    mac.Write(data)
    expectedHMAC := mac.Sum(nil)
    if !hmac.Equal(expectedHMAC, hmacSent) {
        return nil, errors.New("HMAC authentication failed")
    }

    return data, nil
}

// decryptDatagram decrypts the encrypted part of the datagram.
func decryptDatagram(encryptedPart []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    if len(encryptedPart) < aes.BlockSize {
        return nil, errors.New("ciphertext too short")
    }
    iv := encryptedPart[:aes.BlockSize]
    ciphertext := encryptedPart[aes.BlockSize:]
    plaintext := make([]byte, len(ciphertext))
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(plaintext, ciphertext)
    return plaintext, nil
}

// validateAndParseDatagram authenticates and decrypts the datagram,
// populating the provided Datagram pointer with the decrypted data.
func validateAndParseDatagram(buf *[]byte, dg *Datagram) error {
    dg.clientOrServer = (*buf)[0] // Read the ClientOrServer byte

    // Step 1: Populate the Datagram fields
    dg.Username = ToString((*buf)[1:33]) // Populate Username

    // Construct directory path based on the session type
    var dirPath string
    if dg.clientOrServer == 0 { // Client session
        dirPath = filepath.Join(datadir, "accounts", dg.Username)
    } else { // Server session
        dg.PeerUsername = ToString((*buf)[33:65]) // Populate PeerUsername
        dg.PeerServerAddress = ToString((*buf)[65:97]) // Populate PeerServerAddress
        dirPath = filepath.Join(datadir, "accounts", dg.Username, "peers", dg.PeerServerAddress, dg.PeerUsername)
    }

    // Step 2: Load the secret key
    secretKey, err := loadSecretKey(dirPath)
    if err != nil {
        return fmt.Errorf("failed to load secret key: %v", err)
    }

    // Step 3: Authenticate the datagram
    authenticatedData, err := authenticateDatagram(*buf, secretKey)
    if err != nil {
        return fmt.Errorf("failed to authenticate datagram: %v", err)
    }

    // Step 4: Determine the encrypted part based on session type
    var encryptedPart []byte
    if dg.clientOrServer == 0 { // Client session
        encryptedPart = authenticatedData[33:390] // Adjusted for client session encryption range
    } else { // Server session
        encryptedPart = authenticatedData[97:390] // Adjusted for server session encryption range
    }

    // Step 5: Decrypt the datagram
    decryptedData, err := decryptDatagram(encryptedPart, secretKey)
    if err != nil {
        return fmt.Errorf("failed to decrypt datagram: %v", err)
    }

    // Step 6: Write decrypted data back into the Datagram's Arguments field
    if dg.clientOrServer == 0 {
        dg.PeerUsername = ToString((*buf)[33:65])
        dg.PeerServerAddress = ToString((*buf)[65:97])

        peerDir := filepath.Join(datadir, "accounts", dg.Username, "peers", dg.PeerServerAddress, dg.PeerUsername)

        // Inline the peer existence check
        if err := os.Stat(peerDir); err != nil {
            return fmt.Errorf("peer directory does not exist: %v", err)
        }
        copy(dg.Arguments[:], decryptedData[64:]) // Copy the rest to Arguments
    } else {
        // For server sessions, directly copy the decrypted data into Arguments
        copy(dg.Arguments[:], decryptedData)
    }
    // Return nil to indicate success
    return nil
}
