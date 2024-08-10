package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/hmac"
    "crypto/sha256"
    "fmt"
    "os"
    "path/filepath"
)

// loadKeys loads the cryptographic key based on the hash identifier in the datagram
func loadKey(identifier Datagram.Identifier) ([]byte, error) {
    keyDirPath := filepath.Join(datadir, "keys", identifier)

    // Load cryptographic key
    cryptoKey, err = loadSecretKey(keyDirPath, "crypto_key.txt")
    if err != nil {
        return nil, fmt.Errorf("failed to load cryptographic key: %v", err)
    }

    return cryptoKey, nil
}
func decryptPayload(ciphertext, salt, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create AES cipher: %v", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM mode: %v", err)
    }

    plaintext, err := gcm.Open(nil, salt, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("decryption failed: %v", err)
    }

    return plaintext, nil
}
