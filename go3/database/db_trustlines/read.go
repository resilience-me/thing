package db_trustlines

import "ripple/database"

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

// GetSyncCounter retrieves the sync_counter_in value using the datagram to determine the directory.
func GetSyncCounter(dg *Datagram) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.GetUint32FromFile(trustlineDir, "sync_counter.txt")
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
