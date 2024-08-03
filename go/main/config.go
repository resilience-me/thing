package main

import (
    "fmt"
    "os"
    "path/filepath"
)

var serverAddress []byte // Store the server address as a byte array

// GetServerAddress returns the server address as a byte slice
func GetServerAddress() []byte {
    return serverAddress
}

// loadServerAddress reads the server address from the configuration file.
func loadServerAddress() error {
    addressPath := filepath.Join(datadir, "server_address.txt")
    address, err := os.ReadFile(addressPath)
    if err != nil {
        return fmt.Errorf("error loading server address: %w", err)
    }
    serverAddress = address
    return nil
}

// InitConfig initializes the configuration
func initConfig() error {
    if err := loadServerAddress(); err != nil {
        return err
    }
    return nil
}
