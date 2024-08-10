package main

import (
    "crypto/aes"
    "crypto/cipher"
    "fmt"
    "os"
    "path/filepath"
)

type DatagramParser struct {
    Identifier string
    Salt       []byte
    Ciphertext []byte
}

// newDatagramParser creates a new instance of DatagramParser from a buffer
func newDatagramParser(buf []byte) DatagramParser {
    parser := DatagramParser{
        Identifier: string(buf[:32]),   // Convert the first 32 bytes to a string for Identifier
        Salt:       buf[32:44],         // Extract the next 12 bytes as Salt
        Ciphertext: buf[44:],           // The rest is the Ciphertext
    }

    return parser
}

// parseTransaction parses the decrypted plaintext into a Transaction struct
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

// loadKey loads the cryptographic key based on the hash identifier in the datagram
func (dp *DatagramParser) loadKey() ([]byte, error) {
    keyDirPath := filepath.Join(datadir, "keys", dp.Identifier)

    // Load cryptographic key
    cryptoKey, err := loadSecretKey(keyDirPath, "crypto_key.txt")
    if err != nil {
        return nil, fmt.Errorf("failed to load cryptographic key: %v", err)
    }

    return cryptoKey, nil
}

// decryptPayload decrypts the ciphertext using AES-GCM
func (dp *DatagramParser) decryptPayload(key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create AES cipher: %v", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM mode: %v", err)
    }

    plaintext, err := gcm.Open(nil, dp.Salt, dp.Ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("decryption failed: %v", err)
    }

    return plaintext, nil
}

// decryptAndParseDatagram decrypts and parses the datagram into a Transaction
func (dp *DatagramParser) decryptAndParseDatagram() (*Transaction, error) {

    // Load the cryptographic key based on the identifier in the datagram
    secretKey, err := dp.loadKey()
    if err != nil {
        return nil, fmt.Errorf("failed to load cryptographic key: %v", err)
    }

    // Decrypt the payload using AES-GCM
    plaintext, err := dp.decryptPayload(secretKey)
    if err != nil {
        return nil, fmt.Errorf("decryption failed: %v", err)
    }

    // Parse the decrypted payload into the Transaction struct
    tx, err := dp.parseTransaction(plaintext)
    if err != nil {
        return nil, fmt.Errorf("failed to parse transaction: %v", err)
    }

    return tx, nil
}
