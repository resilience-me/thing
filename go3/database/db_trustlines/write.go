package db_trustlines

import (
	"ripple/database"
)

// SetTrustlineOut sets the outbound trustline amount.
func SetTrustlineOut(dg *Datagram, value uint32) error {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.WriteUint32ToFile(trustlineDir, "trustline_out.txt", value)
}

// SetTrustlineOut sets the inbound trustline amount.
func SetTrustlineIn(dg *Datagram, value uint32) error {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.WriteUint32ToFile(trustlineDir, "trustline_in.txt", value)
}

// SetCounter sets the counter value.
func SetCounter(dg *Datagram, value uint32) error {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.WriteUint32ToFile(trustlineDir, "counter.txt", value)
}

// SetCounterOut sets the counter_out value.
func SetCounterOut(dg *Datagram, value uint32) error {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.WriteUint32ToFile(trustlineDir, "counter_out.txt", value)
}

// SetCounterIn sets the counter_in value.
func SetCounterIn(dg *Datagram, value uint32) error {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.WriteUint32ToFile(trustlineDir, "counter_in.txt", value)
}

// SetSyncCounter sets the sync_counter value.
func SetSyncCounter(dg *Datagram, value uint32) error {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.WriteUint32ToFile(trustlineDir, "sync_counter.txt", value)
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

// SetTimestamp sets the sync timestamp.
func SetTimestamp(dg *Datagram, timestamp int64) error {
	trustlineDir := database.GetTrustlineDir(dg)
	return database.WriteTimeToFile(trustlineDir, "timestamp.txt", timestamp)
}
