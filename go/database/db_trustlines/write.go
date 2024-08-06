package db_trustlines

import (
	"resilience/database"
)

// SetTrustlineOut sets the outbound trustline amount.
func SetTrustlineOut(dg *Datagram, value uint32) error {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.WriteUint32ToFile(trustlineDir, "trustline_out.txt", value)
}

// SetCounter sets the counter value.
func SetCounter(dg *Datagram, value uint32) error {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.WriteUint32ToFile(trustlineDir, "counter.txt", value)
}

// SetSyncIn sets the sync_in value.
func SetSyncIn(dg *Datagram, value uint32) error {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.WriteUint32ToFile(trustlineDir, "sync_in.txt", value)
}

// SetSyncOut sets the sync_out value.
func SetSyncOut(dg *Datagram, value uint32) error {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.WriteUint32ToFile(trustlineDir, "sync_out.txt", value)
}

// SetSyncCounterOut sets the sync_counter_out value.
func SetSyncCounterOut(dg *Datagram, value uint32) error {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.WriteUint32ToFile(trustlineDir, "sync_counter_out.txt", value)
}

// SetSyncCounterIn sets the sync_counter_in value.
func SetSyncCounterIn(dg *Datagram, value uint32) error {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.WriteUint32ToFile(trustlineDir, "sync_counter_in.txt", value)
}

// SetTimestamp sets the sync timestamp.
func SetTimestamp(dg *Datagram, timestamp int64) error {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.WriteTimeToFile(trustlineDir, "timestamp.txt", timestamp)
}
