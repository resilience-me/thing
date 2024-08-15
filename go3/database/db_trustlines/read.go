package db_trustlines

import (
	"ripple/main"
	"ripple/database"
)

// GetTrustlineOut retrieves the outbound trustline using the datagram to determine the directory.
func GetTrustlineOut(dg *main.Datagram) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.GetUint32FromFile(trustlineDir, "trustline_out.txt")
}

// GetTrustlineIn retrieves the inbound trustline using the datagram to determine the directory.
func GetTrustlineIn(dg *main.Datagram) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.GetUint32FromFile(trustlineDir, "trustline_in.txt")
}

// GetSyncCounter retrieves the sync_counter_in value using the datagram to determine the directory.
func GetSyncCounter(dg *main.Datagram) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.GetUint32FromFile(trustlineDir, "sync_counter.txt")
}

// GetSyncIn retrieves the sync_in value using the datagram to determine the directory.
func GetSyncIn(dg *main.Datagram) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.GetUint32FromFile(trustlineDir, "sync_in.txt")
}

// GetSyncOut retrieves the sync_out value using the datagram to determine the directory.
func GetSyncOut(dg *main.Datagram) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.GetUint32FromFile(trustlineDir, "sync_out.txt")
}
