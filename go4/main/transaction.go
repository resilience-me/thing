// transaction.go
package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
)

const (
    OffsetTransactionNumber = 0
    OffsetValidator         = OffsetTransactionNumber + 4
    OffsetFromUsername      = OffsetValidator + 32
    OffsetFromServerAddress  = OffsetFromUsername + 32
    OffsetToUsername        = OffsetFromServerAddress + 64
    OffsetToServerAddress   = OffsetToUsername + 32
    OffsetData              = OffsetToServerAddress + 64
    OffsetPreviousHash      = OffsetData + 256
    OffsetSignature         = OffsetPreviousHash + 32
)

// Transaction struct definition
type Transaction struct {
    TransactionNumber [4]byte // Optional transaction number to track order
    Validator         [32]byte // Public key or identifier of the validator
    FromUsername      [32]byte // Initiating user's identifier
    FromServerAddress [64]byte // Address of the initiating user's server
    ToUsername        [32]byte // Receiving user's identifier
    ToServerAddress   [64]byte // Address of the receiving user's server
    Data              [256]byte // The first byte here represents the Command
    PreviousHash      [32]byte // Hash of the previous transaction in the chain
    Signature         [64]byte // Digital signature using the validator's private key
}

// fetchLastTransaction retrieves the raw bytes of the last transaction from the file.
func fetchLastTransaction(filename string) ([]byte, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Move to the end of the file
    stat, err := file.Stat()
    if err != nil {
        return nil, err
    }

    // Calculate the size of the transaction struct
    txSize := 4 + 32 + 32 + 32 + 256 + 32 + 64 // Total size of Transaction struct
    offset := stat.Size() - txSize

    // Read the raw bytes of the last transaction
    data := make([]byte, txSize)
    _, err = file.ReadAt(data, offset)
    if err != nil {
        return nil, err
    }

    return data, nil
}

// fetchLastTransactionHash computes the hash of the last transaction.
func fetchLastTransactionHash(filename string) ([32]byte, error) {
    data, err := fetchLastTransaction(filename)
    if err != nil {
        return [32]byte{}, err
    }

    // Compute the hash of the raw bytes
    hash := sha256.Sum256(data)

    return hash, nil
}

// AppendTransaction appends a transaction to the file
func AppendTransaction(filePath string, tx Transaction) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data := tx.Serialize()
	_, err = file.Write(data)
	return err
}

// ReadTransactions reads all transactions from the file
func ReadTransactions(filePath string) ([]Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var transactions []Transaction
	for {
		var tx Transaction
		err = binary.Read(file, binary.LittleEndian, &tx)
		if err != nil {
			break
		}
		transactions = append(transactions, tx)
	}

	if err.Error() != "EOF" {
		return nil, err
	}

	return transactions, nil
}
