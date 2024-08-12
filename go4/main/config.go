package main

import (
    "os"
    "path/filepath"
)

var datadir = filepath.Join(os.Getenv("HOME"), "ripple")

// GetDataDir returns the datadir as a string
func GetDataDir() string {
    return datadir
}
