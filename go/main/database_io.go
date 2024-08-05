package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// getUint32FromFile reads the contents of a file, parses it as a uint32, and returns the value.
func getUint32FromFile(dir, filename string) (uint32, error) {
	filePath := filepath.Join(dir, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	// Convert the file content to uint32
	value, err := strconv.ParseUint(string(data), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("error parsing value from file %s: %v", filePath, err)
	}
	return uint32(value), nil
}

// GetCounterOut retrieves the outbound counter using the datagram to determine the directory.
func GetCounterOut(dg *Datagram) (uint32, error) {
	trustlineDir := GetTrustlineDir(dg)
	return getUint32FromFile(trustlineDir, "counter_out.txt")
}

// GetCounterIn retrieves the inbound counter using the datagram to determine the directory.
func GetCounterIn(dg *Datagram) (uint32, error) {
	trustlineDir := GetTrustlineDir(dg)
	return getUint32FromFile(trustlineDir, "counter_in.txt")
}

// GetTrustlineOut retrieves the outbound trustline using the datagram to determine the directory.
func GetTrustlineOut(dg *Datagram) (uint32, error) {
	trustlineDir := GetTrustlineDir(dg)
	return getUint32FromFile(trustlineDir, "trustline_out.txt")
}

// GetTrustlineIn retrieves the inbound trustline using the datagram to determine the directory.
func GetTrustlineIn(dg *Datagram) (uint32, error) {
	trustlineDir := GetTrustlineDir(dg)
	return getUint32FromFile(trustlineDir, "trustline_in.txt")
}
