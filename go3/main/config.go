package main

import (
    "io"
    "log"
    "os"
    "path/filepath"
)

var serverAddress string // Store the server address as a string

// GetServerAddress returns the server address as a string
func GetServerAddress() string {
    return serverAddress
}

// loadServerAddress reads the server address from the configuration file.
func loadServerAddress() error {
    addressPath := filepath.Join(datadir, "server_address.txt")
    address, err := os.ReadFile(addressPath)
    if err != nil {
        // Log the error with details rather than returning it with fmt.Errorf
        log.Printf("Error loading server address from %s: %v", addressPath, err)
        return err
    }
    serverAddress = string(address)
    log.Printf("Loaded server address: %s", serverAddress) // Log that the address was loaded
    return nil
}

func setupLogger() {
    // Create or open a log file
    logFile, err := os.OpenFile("ripple.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err) // log.Fatalf logs to stderr and exits the program
    }
    log.SetOutput(logFile)
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// InitConfig initializes the configuration
func initConfig() error {
    setupLogger()
    log.Println("Logger setup completed, initializing configuration...")

    // Load server address, handle its error
    if err := loadServerAddress(); err != nil {
        log.Printf("Failed to load server address: %v", err)
        return err
    }

    log.Println("Configuration initialized successfully.")
    return nil
}
