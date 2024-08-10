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

type DatagramParser struct {
    Identifier []byte
    Salt       []byte
    Ciphertext []byte
}

// loadKeys loads the cryptographic key based on the hash identifier in the datagram
func (dp *DatagramParser) loadKey() ([]byte, error) {
    keyDirPath := filepath.Join(datadir, "keys", dp.Identifier)

    // Load cryptographic key
    cryptoKey, err = loadSecretKey(keyDirPath, "crypto_key.txt")
    if err != nil {
        return nil, fmt.Errorf("failed to load cryptographic key: %v", err)
    }

    return cryptoKey, nil
}
func (dp *DatagramParser) decryptPayload(key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create AES cipher: %v", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM mode: %v", err)
    }

    plaintext, err := gcm.Open(nil, dp.salt, dp.ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("decryption failed: %v", err)
    }

    return plaintext, nil
}

func (dp *DatagramParser) parseTransaction(plaintext []byte) (*Transaction, error) {
    tx := &Transaction{
        Command:           plaintext[0],
        Username:          string(plaintext[1:33]),  // Assuming Username is 32 bytes
        PeerUsername:      string(plaintext[33:65]), // Assuming PeerUsername is 32 bytes
        PeerServerAddress: string(plaintext[65:129]), // Assuming PeerServerAddress is 64 bytes
    }
    copy(tx.Arguments[:], plaintext[129:385]) // Copy Arguments (256 bytes)
    copy(tx.Counter[:], plaintext[385:389])   // Copy Counter (4 bytes)

    return tx, nil
}

func (dp *DatagramParser) decryptAndParseDatagram(buf []byte) (*Transaction, error) {
    // Load the cryptographic key based on the identifier in the datagram
    secretKey, err := loadKey(dp.Identifier)
    if err != nil {
        return nil, fmt.Errorf("failed to load cryptographic key: %v", err)
    }

    // Decrypt the payload using AES-GCM
    plaintext, err := decryptPayload(secretKey)
    if err != nil {
        return nil, fmt.Errorf("decryption failed: %v", err)
    }

    // Parse the decrypted payload into the Transaction struct
    tx, err := parseTransaction(plaintext)
    if err != nil {
        return nil, fmt.Errorf("failed to parse transaction: %v", err)
    }

    return tx, nil
}
