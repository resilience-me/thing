// transaction.go
package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
)

type Transaction struct {
    TransactionNumber [4]byte // Optional transaction number to track order
    Validator         [32]byte // Public key or identifier of the validator
    From              [32]byte // Initiating user's identifier
    To                [32]byte // Receiving user's identifier
    Data              [256]byte // The first byte here represents the Command
    PreviousHash      [32]byte // Hash of the previous transaction in the chain
    Signature         [64]byte // Digital signature using the validator's private key
}

// CalculateHash computes the SHA-256 hash of a transaction
func CalculateHash(tx Transaction) [32]byte {
	// Create a new SHA-256 hash
	hash := sha256.New()

	// Create a buffer to hold the transaction bytes
	var buffer bytes.Buffer

	// Write the transaction data to the buffer
	binary.Write(&buffer, binary.LittleEndian, tx)

	// Write the buffer content to the hash
	hash.Write(buffer.Bytes())

	// Return the computed hash
	var hashArray [32]byte
	copy(hashArray[:], hash.Sum(nil))
	return hashArray
}


// AppendTransaction appends a new transaction to the specified file
func AppendTransaction(filePath string, tx Transaction) error {
	// Open the file in append mode, create if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the transaction to the file
	_, err = file.Write(tx[:]) // Directly write the entire struct as bytes
	if err != nil {
		return err
	}

	return nil
}

// ReadTransactions reads all transactions from the specified file
func ReadTransactions(filePath string) ([]Transaction, error) {
	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the entire file content
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	numTransactions := fileSize / int64(binary.Size(Transaction{})) // Calculate the number of transactions

	transactions := make([]Transaction, numTransactions)

	// Read each transaction
	for i := 0; i < int(numTransactions); i++ {
		var tx Transaction
		err := binary.Read(file, binary.LittleEndian, &tx)
		if err != nil {
			return nil, err
		}
		transactions[i] = tx
	}

	return transactions, nil
}

// FetchLastTransaction reads the last transaction from the specified file
func FetchLastTransaction(filePath string) (Transaction, error) {
	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return Transaction{}, err
	}
	defer file.Close()

	// Seek to the end of the file
	if _, err := file.Seek(-int64(binary.Size(Transaction{})), os.SEEK_END); err != nil {
		return Transaction{}, err
	}

	// Read the last transaction
	var lastTx Transaction
	err = binary.Read(file, binary.LittleEndian, &lastTx)
	if err != nil {
		return Transaction{}, err
	}

	return lastTx, nil
}
