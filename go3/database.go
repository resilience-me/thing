package main

import (
    "os"
    "path/filepath"
)

var datadir = filepath.Join(os.Getenv("HOME"), "ripple")

// checkDirExists checks if a specific directory exists.
func checkDirExists(dirPath string) (bool, error) {
    // Use os.Stat to get the file info
    _, err := os.Stat(dirPath)
    if err != nil {
        if os.IsNotExist(err) {
            // Directory does not exist
            return false, nil
        }
        // Some other error occurred during Stat
        return false, err
    }

    // Directory exists
    return true, nil
}
