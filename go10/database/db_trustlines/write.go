package db_trustlines

import (
	"ripple/types"
	"ripple/database"
)

// SetTrustlineOut sets the outbound trustline amount.
func SetTrustlineOut(username, peerServerAddress, peerUsername string, value uint32) error {
	trustlineDir := database.GetTrustlineDir(username, peerServerAddress, peerUsername)
	return database.WriteUint32ToFile(trustlineDir, "trustline_out.txt", value)
}

// SetTrustlineOut sets the inbound trustline amount.
func SetTrustlineIn(username, peerServerAddress, peerUsername string, value uint32) error {
	trustlineDir := database.GetTrustlineDir(username, peerServerAddress, peerUsername)
	return database.WriteUint32ToFile(trustlineDir, "trustline_in.txt", value)
}

// SetSyncCounter sets the sync_counter value.
func SetSyncCounter(dg *types.Datagram, value uint32) error {
	trustlineDir := database.GetTrustlineDir(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
	return database.WriteUint32ToFile(trustlineDir, "sync_counter.txt", value)
}

// SetSyncIn sets the sync_in value.
func SetSyncIn(dg *types.Datagram, value uint32) error {
	trustlineDir := database.GetTrustlineDir(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
	return database.WriteUint32ToFile(trustlineDir, "sync_in.txt", value)
}

// SetSyncOut sets the sync_out value.
func SetSyncOut(dg *types.Datagram, value uint32) error {
	trustlineDir := database.GetTrustlineDir(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
	return database.WriteUint32ToFile(trustlineDir, "sync_out.txt", value)
}

// SetTimestamp sets the sync timestamp.
func SetTimestamp(dg *types.Datagram, timestamp int64) error {
	trustlineDir := database.GetTrustlineDir(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
	return database.WriteTimeToFile(trustlineDir, "timestamp.txt", timestamp)
}
