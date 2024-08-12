package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// LoadSharedKey loads the shared symmetric key from a file based on the given identifier.
func LoadSharedKey(datadir string, identifier [20]byte) ([]byte, error) {
	// Convert identifier to a string representation (hex)
	identifierStr := fmt.Sprintf("%x", identifier)

	// Construct the file path in datadir/keys/ with the identifier as the filename and .key extension
	filePath := filepath.Join(datadir, "keys", identifierStr+".key")

	// Read the key from the file
	key, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load shared key from file %s: %v", filePath, err)
	}

	return key, nil
}

// SaveSharedKey saves the shared symmetric key to a file based on the given identifier.
func SaveSharedKey(datadir string, identifier [20]byte, key []byte) error {
	// Convert identifier to a string representation (hex)
	identifierStr := fmt.Sprintf("%x", identifier)

	// Construct the file path in datadir/keys/ with the identifier as the filename and .key extension
	filePath := filepath.Join(datadir, "keys", identifierStr+".key")

	// Write the key to the file
	err := ioutil.WriteFile(filePath, key, 0600) // 0600 ensures that only the owner can read/write the file
	if err != nil {
		return fmt.Errorf("failed to save shared key to file %s: %v", filePath, err)
	}

	return nil
}
