package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
)

var serverAddress string // Store the server address as a byte array

// GetServerAddress returns the server address as a byte slice
func GetServerAddress() string {
    return serverAddress
}

// loadServerAddress reads the server address from the configuration file.
func loadServerAddress() error {
    addressPath := filepath.Join(datadir, "server_address.txt")
    address, err := os.ReadFile(addressPath)
    if err != nil {
        return fmt.Errorf("error loading server address: %w", err)
    }
    serverAddress = string(address)
    return nil
}

func setupLogger() {
    // Create or open a log file
    logFile, err := os.OpenFile("ripple.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err) // log.Fatalf logs to stderr and calls os.Exit(1)
    }
    // Set up logging to both the file and stderr
    multiWriter := io.MultiWriter(os.Stdout, logFile)
    log.SetOutput(multiWriter)

    // Optional: Set the logging format to include the date, time, and file source
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// InitConfig initializes the configuration
func initConfig() error {
    setupLogger()
    // Log that the logger has been successfully set up
    log.Println("Logger setup completed, initializing configuration...")

    // Example function to load server address, handle its error
    if err := loadServerAddress(); err != nil {
        log.Printf("Failed to load server address: %v", err)
        return err
    }

    log.Println("Configuration initialized successfully.")
    return nil
}
