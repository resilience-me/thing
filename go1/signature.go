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
func loadSecretKeyFromDir(dir string) ([]byte, error) {
    secretKeyPath := filepath.Join(dir, "secretkey.txt")
    secretKey, err := os.ReadFile(secretKeyPath)
    if err != nil {
        return nil, fmt.Errorf("error reading secret key from %s: %w", secretKeyPath, err)
    }

    return secretKey, nil
}

func loadSecretKey(buf []byte) ([]byte, error) {
    clientOrServer := buf[0]

    var dirPath string
    if clientOrServer == 0 { // Client session
        username := ToString(buf[1:33]) // Convert [32]byte to a slice and trim
        dirPath = filepath.Join(datadir, "accounts", username)
    } else { // Server session
        username := ToString(buf[1:33]) // Convert [32]byte to a slice and trim
        peerUsername := ToString(buf[33:65]) // Convert [32]byte to a slice and trim
        peerServerAddress := ToString(buf[65:97]) // Convert [32]byte to a slice and trim
        dirPath = filepath.Join(datadir, "accounts", username, "peers", peerServerAddress, peerUsername)
    }

    // Load the secret key from the constructed directory path
    secretKey, err := loadSecretKeyFromDir(dirPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load secret key: %v", err)
    }
    return secretKey, nil
}

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

func authenticateAndDecrypt(buf *[]byte) ([]byte, error) {
    // Step 1: Load the secret key
    secretKey, err := loadSecretKey(*buf)
    if err != nil {
        return nil, fmt.Errorf("failed to load secret key: %v", err)
    }

    // Step 2: Authenticate the datagram
    authenticatedData, err := authenticateDatagram(*buf, secretKey)
    if err != nil {
        return nil, fmt.Errorf("failed to authenticate datagram: %v", err)
    }

    // Step 3: Determine the encrypted part based on session type
    var encryptedPart []byte
    if (*buf)[0] == 0 { // Client session
        encryptedPart = authenticatedData[33:390] // Adjusted for client session encryption range
    } else { // Server session
        encryptedPart = authenticatedData[97:390] // Adjusted for server session encryption range
    }

    // Step 4: Decrypt the datagram
    decryptedData, err := decryptDatagram(encryptedPart, secretKey)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt datagram: %v", err)
    }

    // Return the decrypted data
    return decryptedData, nil
}
