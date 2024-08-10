type Datagram struct {
    Identifier []byte
    Salt       []byte
    Ciphertext []byte
}

func parseTransaction(plaintext []byte) (*Transaction, error) {

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

func decryptAndParseDatagram(dg *Datagram) (*Transaction, error) {
    // Load the cryptographic key based on the hash identifier in the buffer
    cryptoKey, err := loadKey(dg.Identifier)
    if err != nil {
        return nil, fmt.Errorf("failed to load cryptographic key: %v", err)
    }

    // Decrypt the payload using AES-GCM
    plaintext, err := decryptPayload(dg.Ciphertext, cryptoKey)
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

