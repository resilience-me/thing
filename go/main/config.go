package main

import (
    "fmt"
    "os"
    "path/filepath"
)

// loadServerAddress reads the server address from server_address.txt
func loadServerAddress() (string, error) {
    addressPath := filepath.Join(datadir, "server_address.txt")
    address, err := os.ReadFile(addressPath)
    if err != nil {
        return "", fmt.Errorf("error reading server address file: %w", err)
    }

    return string(address), nil
}
