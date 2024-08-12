// GenerateSharedKey generates a shared symmetric key using ECDH key exchange.
func GenerateSharedKey(privKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey) ([]byte, error) {
    // Perform ECDH key exchange
    x, _ := pubKey.Curve.ScalarMult(pubKey.X, pubKey.Y, privKey.D.Bytes())

    // Derive a symmetric key from the shared secret using SHA-256
    sharedKey := sha256.Sum256(x.Bytes())
    return sharedKey[:], nil
}

// EncryptTransactionRequest encrypts the signed transaction request with the shared symmetric key.
func EncryptTransactionRequest(request *TransactionRequest, sharedKey []byte) ([]byte, error) {
    // Serialize the signed transaction request
    data := append(request.From[:], request.To[:]...)
    data = append(data, request.Data[:]...)
    data = append(data, request.Signature...)

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

    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    return ciphertext, nil
}

// DecryptTransactionRequest decrypts the transaction request with the shared symmetric key.
func DecryptTransactionRequest(ciphertext, sharedKey []byte) (*TransactionRequest, error) {
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

    // Deserialize the decrypted data into a TransactionRequest
    var request TransactionRequest
    copy(request.From[:], decryptedData[:32])
    copy(request.To[:], decryptedData[32:64])
    copy(request.Data[:], decryptedData[64:320])
    request.Signature = decryptedData[320:]

    return &request, nil
}

// VerifyTransactionRequest verifies the signature of the transaction request.
func VerifyTransactionRequest(request *TransactionRequest, pubKey *ecdsa.PublicKey) bool {
    // Serialize the request data to be verified
    dataToVerify := append(request.From[:], request.To[:]...)
    dataToVerify = append(dataToVerify, request.Data[:]...)

    // Hash the data
    hash := sha256.Sum256(dataToVerify)

    // Extract r and s from the signature
    r := new(big.Int).SetBytes(request.Signature[:32])
    s := new(big.Int).SetBytes(request.Signature[32:])

    // Verify the signature
    return ecdsa.Verify(pubKey, hash[:], r, s)
}
