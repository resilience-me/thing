package db_trustlines

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"ripple/database"
)

// GetTrustlineOut retrieves the outbound trustline using the datagram to determine the directory.
func GetTrustlineOut(dg *Datagram) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.GetUint32FromFile(trustlineDir, "trustline_out.txt")
}

// GetTrustlineIn retrieves the inbound trustline using the datagram to determine the directory.
func GetTrustlineIn(dg *Datagram) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.GetUint32FromFile(trustlineDir, "trustline_in.txt")
}

// GetCounter retrieves the counter value using the datagram to determine the directory.
func GetCounter(dg *Datagram) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.GetUint32FromFile(trustlineDir, "counter.txt")
}

// GetCounterOut retrieves the sync_counter_out value using the datagram to determine the directory.
func GetSyncCounterOut(dg *Datagram) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.GetUint32FromFile(trustlineDir, "counter_out.txt")
}

// GetSyncCounterIn retrieves the sync_counter_in value using the datagram to determine the directory.
func GetCounterIn(dg *Datagram) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.GetUint32FromFile(trustlineDir, "counter_in.txt")
}

// GetSyncIn retrieves the sync_in value using the datagram to determine the directory.
func GetSyncIn(dg *Datagram) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.GetUint32FromFile(trustlineDir, "sync_in.txt")
}

// GetSyncOut retrieves the sync_out value using the datagram to determine the directory.
func GetSyncOut(dg *Datagram) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.GetUint32FromFile(trustlineDir, "sync_out.txt")
}
