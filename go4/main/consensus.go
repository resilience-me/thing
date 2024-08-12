package main

import (
    "fmt"
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

// CreateDatagram encrypts the TransactionRequest, generates the identifier, and creates the Datagram.
func CreateDatagram(txRequest TransactionRequest, sharedKey []byte) (Datagram, error) {
	// Serialize the TransactionRequest
	txRequestBytes := append(txRequest.From[:], append(txRequest.To[:], append(txRequest.Data[:], txRequest.Signature[:]...)...)...)

	// Encrypt the TransactionRequest
	encryptedRequest, err := EncryptTransactionRequest(txRequestBytes, sharedKey)
	if err != nil {
		return Datagram{}, fmt.Errorf("error encrypting transaction request: %v", err)
	}

	// Generate the identifier based on the From and To addresses
	identifier := generateIdentifier(txRequest.From, txRequest.To)

	// Create the Datagram
	var datagram Datagram
	copy(datagram.Identifier[:], identifier[:])
	copy(datagram.Payload[:], encryptedRequest[:])

	return datagram, nil
}
