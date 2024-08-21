package auth

import (
    "crypto/hmac"
    "crypto/sha256"
    "fmt"
    "os"
    "path/filepath"

    "ripple/database"
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

func loadClientSecretKey(dg *types.Datagram) ([]byte, error) {
    accountDir := database.GetAccountDir(dg)
    return loadSecretKeyFromDir(accountDir)
}

func LoadServerSecretKey(dg *types.Datagram) ([]byte, error) {
    peerDir := database.GetPeerDir(dg)
    return loadSecretKeyFromDir(peerDir)
}

func LoadServerSecretKeyOut(dg *types.Datagram, peerServerAddress string) ([]byte, error) {
    peerDir := database.GetPeerDirOut(dg, peerServerAddress)
    return loadSecretKeyFromDir(peerDir)
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
func GenerateHMAC(data []byte, secret []byte) ([]byte, error) {
    h := hmac.New(sha256.New, secret)
    _, err := h.Write(data)
    if err != nil {
        return nil, fmt.Errorf("failed to write data to HMAC: %w", err)
    }
    signature := h.Sum(nil) // Get the raw byte slice of the HMAC
    return signature, nil    // Return the signature as a byte slice
}
