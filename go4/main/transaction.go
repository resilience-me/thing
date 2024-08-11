// transaction.go
package main

import (
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

