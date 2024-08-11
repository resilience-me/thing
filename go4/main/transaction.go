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
    OffsetFromServerAddress = OffsetFromUsername + 32
    OffsetToUsername        = OffsetFromServerAddress + 32
    OffsetToServerAddress   = OffsetToUsername + 32
    OffsetData              = OffsetToServerAddress + 32
    OffsetPreviousHash      = OffsetData + 256
    OffsetSignature         = OffsetPreviousHash + 32
    LengthTransaction       = OffsetSignature + 32
)

// Transaction struct definition
type Transaction struct {
    TransactionNumber [4]byte // Optional transaction number to track order
    Validator         [32]byte // Public key or identifier of the validator
    FromUsername      [32]byte // Initiating user's identifier
    FromServerAddress [32]byte // Address of the initiating user's server
    ToUsername        [32]byte // Receiving user's identifier
    ToServerAddress   [32]byte // Address of the receiving user's server
    Data              [256]byte // The first byte here represents the Command
    PreviousHash      [32]byte // Hash of the previous transaction in the chain
    Signature         [32]byte // Digital signature using the validator's private key
}
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
)

func signTransaction(privKey *ecdsa.PrivateKey, data []byte) ([]byte, []byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, privKey, data)
	if err != nil {
		return nil, nil, err
	}
	return r.Bytes(), s.Bytes(), nil
}

func verifyTransaction(pubKey *ecdsa.PublicKey, data []byte, rBytes, sBytes []byte) bool {
	r := new(big.Int).SetBytes(rBytes)
	s := new(big.Int).SetBytes(sBytes)
	return ecdsa.Verify(pubKey, data, r, s)
}
