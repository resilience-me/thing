package auth

import (
    "crypto/hmac"
    "crypto/sha256"
    "fmt"
    "os"
    "path/filepath"

    "ripple/database"
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

// verifyHMAC checks the integrity and authenticity of the received buffer
func verifyHMAC(buf []byte, key []byte) bool {
    // The signature is the last 32 bytes of the buffer
    data := buf[:len(buf)-32]
    signature := buf[len(buf)-32:]
    mac := hmac.New(sha256.New, key)
    mac.Write(data)
    expectedMAC := mac.Sum(nil)
    return hmac.Equal(signature, expectedMAC)
}

// GenerateHMAC generates an HMAC signature for the given data using the provided key.
func generateHMAC(data []byte, secret []byte) ([]byte, error) {
    h := hmac.New(sha256.New, secret)
    _, err := h.Write(data)
    if err != nil {
        return nil, fmt.Errorf("failed to write data to HMAC: %w", err)
    }
    signature := h.Sum(nil) // Get the raw byte slice of the HMAC
    return signature, nil    // Return the signature as a byte slice
}
