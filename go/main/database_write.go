package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// writeUint32ToFile writes a uint32 value to a file.
func writeUint32ToFile(dir, filename string, value uint32) error {
	filePath := filepath.Join(dir, filename)
	return os.WriteFile(filePath, []byte(fmt.Sprintf("%d", value)), 0644)
}

// writeTimeToFile writes a Unix timestamp to a file.
func writeTimeToFile(dir, filename string, timestamp int64) error {
	filePath := filepath.Join(dir, filename)
	return os.WriteFile(filePath, []byte(fmt.Sprintf("%d", timestamp)), 0644)
}

// SetTrustlineOut sets the outbound trustline amount.
func SetTrustlineOut(dg *Datagram, value uint32) error {
	trustlineDir := GetTrustlineDir(dg)
	return writeUint32ToFile(trustlineDir, "trustline_out.txt", value)
}

// SetCounter sets the counter value.
func SetCounter(dg *Datagram, value uint32) error {
	trustlineDir := GetTrustlineDir(dg)
	return writeUint32ToFile(trustlineDir, "counter.txt", value)
}

// SetSyncIn sets the sync_in value.
func SetSyncIn(dg *Datagram, value uint32) error {
	trustlineDir := GetTrustlineDir(dg)
	return writeUint32ToFile(trustlineDir, "sync_in.txt", value)
}

// SetSyncOut sets the sync_out value.
func SetSyncOut(dg *Datagram, value uint32) error {
	trustlineDir := GetTrustlineDir(dg)
	return writeUint32ToFile(trustlineDir, "sync_out.txt", value)
}

// SetTimestamp sets the sync timestamp.
func SetTimestamp(dg *Datagram, timestamp int64) error {
	trustlineDir := GetTrustlineDir(dg)
	return writeTimeToFile(trustlineDir, "timestamp.txt", timestamp)
}
