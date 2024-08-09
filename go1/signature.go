func loadSecretKey(buf []byte) ([]byte, error) {
    clientOrServer := buf[0]

    var dirPath string
    if clientOrServer == 0 { // Client session
        username := ToString(buf[1:33]) // Convert [32]byte to a slice and trim
        dirPath = filepath.Join(datadir, "accounts", username)
    } else { // Server session
        username := ToString(buf[1:33]) // Convert [32]byte to a slice and trim
        peerUsername := ToString(buf[33:65]) // Convert [32]byte to a slice and trim
        peerServerAddress := ToString(buf[65:97]) // Convert [32]byte to a slice and trim
        dirPath = filepath.Join(datadir, "accounts", username, "peers", peerServerAddress, peerUsername)
    }

    // Load the secret key from the constructed directory path
    secretKey, err := loadSecretKeyFromDir(dirPath)
    if err != nil {
        return nil, fmt.Errorf("failed to load secret key: %v", err)
    }
    return secretKey, nil
}

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
