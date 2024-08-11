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
    SizeNumber      = 4
    SizeValidator   = 32
    SizeFrom        = 32
    SizeTo          = 32
    SizeData        = 256
    SizeParentHash  = 32
    SizeSignature   = 64

    SizeTransaction = 452

    OffsetNumber    = 0
    OffsetValidator = OffsetNumber + SizeNumber
    OffsetFrom      = OffsetValidator + SizeValidator
    OffsetTo        = OffsetFrom + SizeFrom
    OffsetData      = OffsetTo + SizeTo
    OffsetParentHash = OffsetData + SizeData
    OffsetSignature = OffsetParentHash + SizeParentHash
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

	offset := int64(index * SizeTransaction)
	data := make([]byte, SizeTransaction)
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
	lastTransactionOffset := info.Size() - int64(SizeTransaction)
	transaction := make([]byte, SizeTransaction)
	_, err = file.ReadAt(transaction, lastTransactionOffset)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransactionChainHeight returns the number of transactions stored in the file
func GetTransactionChainHeight(filename string) (uint32, error) {
    file, err := os.Open(filename)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    info, err := file.Stat()
    if err != nil {
        return 0, err
    }

    // Calculate the number of transactions
    transactionCount := uint32(info.Size() / int64(SizeTransaction))
    return transactionCount, nil
}

// ExtractParentHash extracts the ParentHash from the transaction bytes.
func ExtractParentHash(transaction []byte) ([32]byte, error) {
	var parentHash [32]byte
	copy(parentHash[:], transaction[OffsetParentHash:OffsetParentHash+SizeParentHash])
	return parentHash, nil
}
