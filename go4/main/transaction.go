// transaction.go
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
    OffsetNumber 	= 0
    OffsetValidator    	= OffsetNumber + 4
    OffsetFrom      	= OffsetValidator + 32
    OffsetTo        	= OffsetFrom + 32
    OffsetData         	= OffsetTo + 32
    OffsetParentHash    = OffsetData + 256
    OffsetSignature     = OffsetParentHash + 32
    LengthTransaction   = OffsetSignature + 64
)

type Transaction struct {
    Number            [4]byte
    Validator         [32]byte
    From              [32]byte
    To	              [32]byte
    Data              [256]byte
    ParentHash        [32]byte
    Signature         [64]byte
}

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

func writeRawTransactionToFile(data []byte, filename string) error {

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func readRawTransactionFromFile(index int, filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	offset := int64(index * LengthTransaction)
	data := make([]byte, LengthTransaction)
	_, err = file.ReadAt(data, offset)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetLatestTransaction retrieves the raw bytes of the latest transaction in the file.
func GetLatestTransaction(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Calculate the offset for the last transaction
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	lastTransactionOffset := info.Size() - int64(LengthTransaction)
	transaction := make([]byte, LengthTransaction)
	_, err = file.ReadAt(transaction, lastTransactionOffset)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
