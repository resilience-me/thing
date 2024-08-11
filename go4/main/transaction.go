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

// Encode encodes a Transaction into a JSON byte array.
func (t *Transaction) Encode() ([]byte, error) {
    return json.Marshal(t)
}

// Decode decodes a JSON byte array into a Transaction.
func Decode(data []byte) (*Transaction, error) {
    var t Transaction
    err := json.Unmarshal(data, &t)
    if err != nil {
        return nil, err
    }
    return &t, nil
}

// Key generates a unique key for LevelDB based on the TransactionNumber.
func (t *Transaction) Key() []byte {
    key := make([]byte, 4)
    binary.BigEndian.PutUint32(key, binary.BigEndian.Uint32(t.TransactionNumber[:]))
    return key
}
