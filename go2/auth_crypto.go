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
