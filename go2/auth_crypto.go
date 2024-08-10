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
)

// loadKeys loads both cryptographic and authentication keys based on the hash identifier in the buffer
func loadKeys(buf []byte) (cryptoKey, authKey []byte, err error) {
    hashIdentifier := string(buf[:32]) // Assume HashIdentifier is the first 32 bytes
    keyDirPath := filepath.Join(datadir, "keys", hashIdentifier)

    // Load cryptographic key
    cryptoKey, err = loadSecretKey(keyDirPath, "crypto_key.txt")
    if err != nil {
        return nil, nil, fmt.Errorf("failed to load cryptographic key: %v", err)
    }

    // Load authentication key
    authKey, err = loadSecretKey(keyDirPath, "auth_key.txt")
    if err != nil {
        return nil, nil, fmt.Errorf("failed to load authentication key: %v", err)
    }

    return cryptoKey, authKey, nil
}

// authenticatePayload verifies the HMAC from the buffer
func authenticatePayload(buf []byte, authKey []byte) error {
    hmacLength := 32 // Known HMAC length
    payloadLength := len(buf) - hmacLength
    encryptedPayload := buf[:payloadLength]
    receivedHMAC := buf[payloadLength:]

    mac := hmac.New(sha256.New, authKey)
    mac.Write(encryptedPayload) // Authenticate the encrypted data
    expectedMAC := mac.Sum(nil)
    if !hmac.Equal(expectedMAC, receivedHMAC) {
        return fmt.Errorf("HMAC validation failed")
    }
    return nil
}

// decryptPayload decrypts the payload directly into the buffer, removing the hash identifier and HMAC
func decryptPayload(buf []byte, cryptoKey []byte) error {
    iv := buf[:aes.BlockSize]
    ciphertext := buf[aes.BlockSize : len(buf)-32] // Exclude HMAC length

    block, err := aes.NewCipher(cryptoKey)
    if err != nil {
        return nil, fmt.Errorf("failed to create AES cipher: %v", err)
    }
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext) // Decrypt in place

    return ciphertext, nil // Return the decrypted payload
}

// authenticateAndDecrypt authenticates and decrypts the datagram
func authenticateAndDecrypt(buf []byte) error {
    // Load both cryptographic and authentication keys
    cryptoKey, authKey, err := loadKeys(buf)
    if err != nil {
        return err // error already formatted
    }

    // Authenticate the encrypted payload first
    if err := authenticatePayload(buf, authKey); err != nil {
        return fmt.Errorf("authentication failed: %v", err)
    }

    // Decrypt the payload after authentication is successful
    if err := decryptPayload(buf, cryptoKey); err != nil {
        return fmt.Errorf("decryption failed: %v", err)
    }

    return nil
}
