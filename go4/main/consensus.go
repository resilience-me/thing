package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/ecdsa"
    "crypto/rand"
    "crypto/sha256"
    "encoding/binary"
    "fmt"
    "io"
)

// HandleTransactionRequest processes a transaction request from the non-validator.
func HandleTransactionRequest(filename string, request []byte, validatorID []byte) error {
    // Ensure that this account is the current validator
    validatedLatestBlock, err := ValidatedLatestBlock(filename, validatorID)
    if err != nil {
        return err
    }
    if validatedLatestBlock {
        return fmt.Errorf("this account is not the current validator")
    }

    // Convert the request into a full transaction
    rawTransaction, err := ConvertRawBytesToTransaction(request)
    if err != nil {
        return fmt.Errorf("failed to convert request to transaction: %v", err)
    }

    copy(rawTransaction[OffsetValidator:], validatorID[:SizeValidator])

    if err := PrepareAndStoreTransaction("transactions.dat", rawTransaction, privateKey); err != nil {
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
