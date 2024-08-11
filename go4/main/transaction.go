// transaction.go
package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
)

// Transaction struct definition
type Transaction struct {
	TransactionNumber [4]byte // Optional transaction number to track order
	Validator         [32]byte // Public key or identifier of the validator
	From              [32]byte // Initiating user's identifier
	To                [32]byte // Receiving user's identifier
	Data              [256]byte // The first byte here represents the Command
	PreviousHash      [32]byte // Hash of the previous transaction in the chain
	Signature         [64]byte // Digital signature using the validator's private key
}

// Serialize serializes the Transaction struct to a byte slice
func (tx *Transaction) Serialize() []byte {
	data := make([]byte, 0)
	data = append(data, tx.TransactionNumber[:]...)
	data = append(data, tx.Validator[:]...)
	data = append(data, tx.From[:]...)
	data = append(data, tx.To[:]...)
	data = append(data, tx.Data[:]...)
	data = append(data, tx.PreviousHash[:]...)
	data = append(data, tx.Signature[:]...)
	return data
}

// CalculateHash computes the SHA-256 hash of a transaction
func CalculateHash(tx Transaction) [32]byte {
	hash := sha256.New()
	hash.Write(tx.Serialize())
	var hashArray [32]byte
	copy(hashArray[:], hash.Sum(nil))
	return hashArray
}

// FetchLastTransaction reads the last transaction from the specified file
func FetchLastTransaction(filePath string) (Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Transaction{}, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return Transaction{}, err
	}

	if stat.Size() < int64(binary.Size(Transaction{})) {
		return Transaction{}, fmt.Errorf("file size is smaller than transaction size")
	}

	if _, err := file.Seek(-int64(binary.Size(Transaction{})), os.SEEK_END); err != nil {
		return Transaction{}, err
	}

	var lastTx Transaction
	err = binary.Read(file, binary.LittleEndian, &lastTx)
	if err != nil {
		return Transaction{}, err
	}

	return lastTx, nil
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
