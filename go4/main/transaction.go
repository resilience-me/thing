// transaction.go
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"bytes"
	"math/big"
)

func VerifyTransaction(pubKey *ecdsa.PublicKey, data []byte, rBytes, sBytes []byte) bool {
	r := new(big.Int).SetBytes(rBytes)
	s := new(big.Int).SetBytes(sBytes)
	return ecdsa.Verify(pubKey, data, r, s)
}

func HashAndVerifyTransaction(pubKey *ecdsa.PublicKey, data, rBytes, sBytes []byte) bool {
	// Hash the transaction data using SHA-256
	hash := sha256.Sum256(data)
	return VerifyTransaction(pubKey, hash[:], rBytes, sBytes)
}

// StripSignatureAndVerify verifies the signature of the transaction.
func StripSignatureAndVerify(rawTransaction []byte, pubKey *ecdsa.PublicKey) bool {
	// Determine the length of the data
	dataLen := len(rawTransaction)

	// Define the starting index for the signature
	signatureStart := dataLen - SizeSignature

	// Extract data excluding the signature
	data := rawTransaction[:signatureStart]

	// Define the size of r and s values
	const rsValues = SizeSignature / 2

	// Extract r and s values from the rawTransaction
	rBytes := rawTransaction[signatureStart : signatureStart+rsValues]
	sBytes := rawTransaction[signatureStart+rsValues : dataLen]

	// Verify the signature
	return HashAndVerifyTransaction(pubKey, data, rBytes, sBytes)
}

func SignTransaction(privKey *ecdsa.PrivateKey, data []byte) ([]byte, error) {
    r, s, err := ecdsa.Sign(rand.Reader, privKey, data)
    if err != nil {
        return nil, err
    }
    // Combine r and s into a single byte slice as the signature
    signature := append(r.Bytes(), s.Bytes()...)
    return signature, nil
}

// HashAndSignTransaction hashes the provided data and signs it using the given private key.
// It returns the signature as a byte slice (concatenated r and s values).
func HashAndSignTransaction(privKey *ecdsa.PrivateKey, rawTransaction []byte) ([]byte, error) {
    // Hash the transaction data using SHA-256
    hash := sha256.Sum256(rawTransaction)

    // Sign the hash using ECDSA
    signature, err := SignTransaction(privKey, hash[:])
    if err != nil {
        return nil, fmt.Errorf("failed to sign transaction: %v", err)
    }

    return signature, nil
}

// SignAndInsertSignature hashes the raw transaction data (excluding the signature),
// signs it, and then inserts the signature into the correct position within the raw transaction.
func SignAndInsertSignature(rawTransaction []byte, privKey *ecdsa.PrivateKey) error {
    // Determine the length of the raw transaction data excluding the signature
    signatureOffset := len(rawTransaction) - SizeSignature

    // Hash and sign the transaction data excluding the signature field
    signature, err := HashAndSignTransaction(privKey, rawTransaction[:signatureOffset])
    if err != nil {
        return err
    }

    // Copy the signature into the raw transaction at the correct offset
    copy(rawTransaction[signatureOffset:], signature)

    return nil
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

// GetLatestTransaction retrieves the raw bytes of the latest transaction in the file.
func GetLatestTransaction(filename string) ([]byte, error) {
    // Get the height of the transaction chain
    chainHeight, err := GetTransactionChainHeight(filename)
    if err != nil {
        return nil, err
    }

    // Retrieve the latest transaction
    latestTransaction, err := readRawTransactionFromFile(chainHeight-1, filename)
    if err != nil {
        return nil, err
    }

    return latestTransaction, nil
}

// ExtractParentHash extracts the ParentHash from the transaction bytes and returns it as a []byte.
func ExtractParentHash(transaction []byte) []byte {
    return transaction[OffsetParentHash : OffsetParentHash+SizeParentHash]
}

// FetchLatestValidator retrieves the Validator from the most recent transaction in the file.
func FetchLatestValidator(filename string) ([]byte, error) {
    // Retrieve the latest transaction using the GetLatestTransaction wrapper
    latestTransaction, err := GetLatestTransaction(filename)
    if err != nil {
        return nil, err
    }

    // Extract the Validator from the transaction bytes
    validator := latestTransaction[OffsetValidator : OffsetValidator+SizeValidator]

    return validator, nil
}

// ValidatedLatestBlock checks if the given account validated the latest block
func ValidatedLatestBlock(filename string, accountID []byte) (bool, error) {
    // Fetch the most recent validator
    latestValidator, err := FetchLatestValidator(filename)
    if err != nil {
        return false, err
    }

    // Compare the latest validator with the account's ID using bytes.Equal
    return bytes.Equal(latestValidator, accountID), nil
}

// PrepareAndStoreTransaction prepares a transaction from raw bytes with Number and ParentHash, signs it, and stores it.
func PrepareAndStoreTransaction(filename string, rawTransaction []byte, privateKey *ecdsa.PrivateKey) error {
    chainHeight, err := GetTransactionChainHeight(filename)
    if err != nil {
        return err
    }

    // Set the transaction number in the raw bytes
    binary.BigEndian.PutUint32(rawTransaction[:SizeNumber], chainHeight)

    // Retrieve the latest transaction to get the ParentHash
    latestTransaction, err := readRawTransactionFromFile(chainHeight-1, filename)
    if err != nil {
        return err
    }

    // Extract and set the ParentHash in the raw bytes
    parentHash := ExtractParentHash(latestTransaction)
    copy(rawTransaction[OffsetParentHash:], parentHash)

    // Sign the transaction data and set the signature
    if err := SignAndInsertSignature(rawTransaction, privateKey); err != nil {
        return fmt.Errorf("failed to sign transaction: %v", err)
    }

    // Store the updated transaction
    return writeRawTransactionToFile(rawTransaction, filename)
}
