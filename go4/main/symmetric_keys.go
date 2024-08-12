package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// getKeyFilePath constructs the file path for a given identifier within the datadir/keys/ directory.
func getKeyFilePath(identifier [20]byte) string {
	// Convert identifier to a string representation (hex)
	identifierStr := fmt.Sprintf("%x", identifier)

	// Construct the file path in datadir/keys/ with the identifier as the filename and .key extension
	return filepath.Join(GetDataDir(), "keys", identifierStr+".key")
}

// LoadSharedKey loads the shared symmetric key from a file based on the given identifier.
func LoadSharedKey(identifier [20]byte) ([]byte, error) {
	// Get the file path using the helper function
	filePath := getKeyFilePath(identifier)

	// Read the key from the file
	key, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load shared key from file %s: %v", filePath, err)
	}

	return key, nil
}

// SaveSharedKey saves the shared symmetric key to a file based on the given identifier.
func SaveSharedKey(identifier [20]byte, key []byte) error {
	// Get the file path using the helper function
	filePath := getKeyFilePath(identifier)

	// Write the key to the file
	err := os.WriteFile(filePath, key, 0600) // 0600 ensures that only the owner can read/write the file
	if err != nil {
		return fmt.Errorf("failed to save shared key to file %s: %v", filePath, err)
	}

	return nil
}
