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

// GetCounter retrieves the counter value using the datagram to determine the directory.
func GetCounter(dg *Datagram) (uint32, error) {
	trustlineDir := GetTrustlineDir(dg)
	return getUint32FromFile(trustlineDir, "counter.txt")
}

// GetSyncIn retrieves the sync_in value using the datagram to determine the directory.
func GetSyncIn(dg *Datagram) (uint32, error) {
	trustlineDir := GetTrustlineDir(dg)
	return getUint32FromFile(trustlineDir, "sync_in.txt")
}

// GetSyncOut retrieves the sync_out value using the datagram to determine the directory.
func GetSyncOut(dg *Datagram) (uint32, error) {
	trustlineDir := GetTrustlineDir(dg)
	return getUint32FromFile(trustlineDir, "sync_out.txt")
}

// GetSyncCounterIn retrieves the sync_counter_in value using the datagram to determine the directory.
func GetSyncCounterIn(dg *Datagram) (uint32, error) {
	trustlineDir := GetTrustlineDir(dg)
	return getUint32FromFile(trustlineDir, "sync_counter_in.txt")
}

// GetSyncCounterOut retrieves the sync_counter_out value using the datagram to determine the directory.
func GetSyncCounterOut(dg *Datagram) (uint32, error) {
	trustlineDir := GetTrustlineDir(dg)
	return getUint32FromFile(trustlineDir, "sync_counter_out.txt")
}
