package config

import (
    "os"
    "path/filepath"
)

const (
    Port = 2012
)

var datadir = filepath.Join(os.Getenv("HOME"), "ripple")
var serverAddress string

// GetServerAddress returns the server address as a string
func GetServerAddress() string {
    return serverAddress
}

// GetDataDir returns the datadir as a string
func GetDataDir() string {
    return datadir
}
