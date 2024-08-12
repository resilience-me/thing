// GenerateSharedKey generates a shared symmetric key using ECDH key exchange.
func GenerateSharedKey(privKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey) ([]byte, error) {
    // Perform ECDH key exchange
    x, _ := pubKey.Curve.ScalarMult(pubKey.X, pubKey.Y, privKey.D.Bytes())

    // Derive a symmetric key from the shared secret using SHA-256
    sharedKey := sha256.Sum256(x.Bytes())
    return sharedKey[:], nil
}

// EncryptTransactionRequest encrypts the signed transaction request with the shared symmetric key.
func EncryptTransactionRequest(request []byte, sharedKey []byte) ([]byte, error) {
    // Encrypt the data using AES-GCM for confidentiality and integrity
    block, err := aes.NewCipher(sharedKey)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %v", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %v", err)
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, fmt.Errorf("failed to generate nonce: %v", err)
    }

    ciphertext := gcm.Seal(nonce, nonce, request, nil)
    return ciphertext, nil
}

// DecryptTransactionRequest decrypts the transaction request with the shared symmetric key.
func DecryptTransactionRequest(ciphertext, sharedKey []byte) ([]byte, error) {
    block, err := aes.NewCipher(sharedKey)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %v", err)
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %v", err)
    }

    nonceSize := gcm.NonceSize()
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

    decryptedData, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt data: %v", err)
    }

    return decryptedData, nil
}

// VerifyTransactionRequest verifies the signature of the transaction request.
func VerifyTransactionRequest(request []byte, pubKey *ecdsa.PublicKey) bool {
    // Hash the relevant data from the request directly
    hash := sha256.Sum256(request[:SizeRequest-SizeSignature]) // Exclude the signature

    // Extract r and s from the signature
    r := new(big.Int).SetBytes(request[SizeRequest-SizeSignature : SizeRequest-SizeSignature+32])
    s := new(big.Int).SetBytes(request[SizeRequest-SizeSignature+32 : SizeRequest])

    // Verify the signature
    return ecdsa.Verify(pubKey, hash[:], r, s)
}
