var datadir = filepath.Join(os.Getenv("HOME"), "ripple")

// loadSecretKey loads the secret key from the specified directory.
func loadSecretKeyFromDir(dir string) ([]byte, error) {
    secretKeyPath := filepath.Join(dir, "secretkey.txt")
    secretKey, err := os.ReadFile(secretKeyPath)
    if err != nil {
        return nil, fmt.Errorf("error reading secret key from %s: %w", secretKeyPath, err)
    }

    return secretKey, nil
}

func loadSecretKey(dg *Datagram) ([]byte, error) {
    var keyDir string
    if dg.Command & 0x80 == 0 {
        keyDir = filepath.Join(datadir, "accounts", username)
    } else {
        keyDir = filepath.Join(datadir, "accounts", username, "peers", peerAddress, peerUsername)
    }
    return loadSecretKeyFromDir(keyDir)
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

func authenticateAndParseDatagram(buf []byte) {
    dg := parseDatagram(buf)
    secretKey := loadSecretKey(&dg)
    return verifyHMAC(buf, secretKey)
}
