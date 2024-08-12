// HandleTransactionRequest processes a transaction request from the non-validator.
func HandleTransactionRequest(filename string, request []byte, validatorID []byte) error {
    // Ensure that this account is the current validator
    ValidatedLatestBlock, err := IsValidator(filename, validatorID)
    if err != nil {
        return err
    }
    if ValidatedLatestBlock {
        return fmt.Errorf("this account is not the current validator")
    }

    // Convert the request into a full transaction
    rawTransaction, err := ConvertRawBytesToTransaction(request)
    if err != nil {
        return fmt.Errorf("failed to convert request to transaction: %v", err)
    }

    copy(rawTransaction[OffsetValidator:], validatorID[:SizeValidator])


    err := PrepareAndStoreTransaction("transactions.dat", rawTransaction, privateKey)
    if err != nil {
        return err
    }

    return nil
}

// ConvertRawBytesToTransaction converts raw bytes of a TransactionRequest to a Transaction by populating the fields.
func ConvertRawBytesToTransaction(request []byte) ([]byte, error) {
    // Create a byte slice to hold the transaction
    rawTransaction := make([]byte, LengthTransaction)

    // Copy the entire TransactionRequest data into the rawTransaction starting at OffsetFrom
    copy(rawTransaction[OffsetFrom:], request[:SizeRequest-SizeSignature])

    return rawTransaction, nil
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
