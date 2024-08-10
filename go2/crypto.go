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
func loadKey(dg Datagram) (cryptoKey []byte, err error) {
    keyDirPath := filepath.Join(datadir, "keys", dg.Identifier)

    // Load cryptographic key
    cryptoKey, err = loadSecretKey(keyDirPath, "crypto_key.txt")
    if err != nil {
        return nil, fmt.Errorf("failed to load cryptographic key: %v", err)
    }

    return cryptoKey, nil
}
func decryptPayload(buf []byte, cryptoKey []byte) ([]byte, error) {
    block, err := aes.NewCipher(cryptoKey)
    if err != nil {
        return nil, fmt.Errorf("failed to create AES cipher: %v", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM mode: %v", err)
    }

    nonceSize := gcm.NonceSize()
    if len(buf) < nonceSize {
        return nil, fmt.Errorf("ciphertext too short")
    }

    nonce, ciphertext := buf[:nonceSize], buf[nonceSize:]

    // Decrypt the ciphertext and verify the authentication tag
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("decryption failed: %v", err)
    }

    return plaintext, nil
}
