package auth

import (
    "crypto/sha256"
    "bytes"
    "ripple/database"
    "ripple/types"
)

func loadClientSecretKey(dg *types.Datagram) ([]byte, error) {
    return database.LoadSecretKey(dg.Username)
}

func loadServerSecretKey(dg *types.Datagram) ([]byte, error) {
    return database.LoadPeerSecretKey(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
}

func loadServerSecretKeyOut(dg *types.Datagram, peerServerAddress string) ([]byte, error) {
    return database.LoadPeerSecretKey(dg.PeerUsername, peerServerAddress, dg.Username)
}

// verifySignature checks the integrity of the received buffer
func verifySignature(buf []byte, key []byte) bool {
    // The signature is the last 32 bytes of the buffer
    data := buf[:len(buf)-32]
    signature := buf[len(buf)-32:]

    // Concatenate data and key
    preimage := append(data, key...)

    // Compute the SHA-256 hash
    hash := sha256.Sum256(preimage)

    // Compare the computed hash with the signature directly using bytes.Equal
    return bytes.Equal(signature, hash[:])
}

// generateSignature generates a SHA-256 hash for the given data using the provided key.
func generateSignature(data []byte, secret []byte) []byte {
    // Concatenate data and secret
    preimage := append(data, secret...)

    // Compute the SHA-256 hash
    hash := sha256.Sum256(preimage)

    // Return the hash as a byte slice
    return hash[:]
}
