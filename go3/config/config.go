package config

import (
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"  // Ensure "time" is imported since PathFindingTimeout uses it
)

// PathFindingTimeout is a global constant that defines the timeout duration for pathfinding operations
const PathFindingTimeout = 5 * time.Minute

// CommitTimeout is a global constant that defines the timeout duration for commits during payment
const CommitTimeout = 10 * time.Minute

var datadir = filepath.Join(os.Getenv("HOME"), "ripple")
var serverAddress string // Store the server address as a string

// GetServerAddress returns the server address as a string
func GetServerAddress() string {
    return serverAddress
}

// GetDataDir returns the datadir as a string
func GetDataDir() string {
    return datadir
}

// loadServerAddress reads the server address from the configuration file.
func loadServerAddress() error {
    addressPath := filepath.Join(datadir, "server_address.txt")
    address, err := os.ReadFile(addressPath)
    if err != nil {
        return fmt.Errorf("error loading server address from %s: %w", addressPath, err)
    }
    serverAddress = string(address)
    log.Printf("Loaded server address: %s", serverAddress) // Log that the address was loaded
    return nil
}

// setupLogger initializes the logging configuration.
func setupLogger() error {
    logFile, err := os.OpenFile("ripple.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
    if err != nil {
        return fmt.Errorf("failed to open log file: %w", err)
    }
    log.SetOutput(logFile)
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    return nil
}

// initConfig initializes the configuration by setting up the logger and loading the server address.
func initConfig() error {
    if err := setupLogger(); err != nil {
        return fmt.Errorf("initializing logger: %w", err)
    }
    log.Println("Logger setup completed, initializing configuration...")

    if err := loadServerAddress(); err != nil {
        return fmt.Errorf("initializing configuration by loading server address: %w", err)
    }

    log.Println("Configuration initialized successfully.")
    return nil
}
