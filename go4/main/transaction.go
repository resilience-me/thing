// transaction.go
package main

import (
    "encoding/binary"
    "encoding/json"
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

// AppendTransaction appends a new transaction to the JSON file
func appendTransaction(accountUsername, peerUsername string, transaction Transaction) error {
    relationshipFilePath := filepath.Join("data", "accounts", accountUsername, "peers", peerUsername, "relationship.json")

    // Open the file for appending
    file, err := os.OpenFile(relationshipFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    // Not the first entry; append a comma for JSON array
    _, err := file.WriteString(",")
    if err != nil {
        return err
    }

    // Marshal the new transaction to JSON
    transactionJSON, err := json.Marshal(transaction)
    if err != nil {
        return err
    }

    // Write the new transaction to the file
    _, err = file.Write(transactionJSON)
    if err != nil {
        return err
    }

    return nil
}
