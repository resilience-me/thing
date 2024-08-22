package db_trustlines

import (
	"ripple/types"
	"ripple/database"
)

// GetTrustlineOut retrieves the outbound trustline
func GetTrustlineOut(username, peerServerAddress, peerUsername string) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(username, peerServerAddress, peerUsername)
	return database.GetUint32FromFile(trustlineDir, "trustline_out.txt")
}

// GetTrustlineIn retrieves the inbound trustline
func GetTrustlineIn(username, peerServerAddress, peerUsername string) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(username, peerServerAddress, peerUsername)
	return database.GetUint32FromFile(trustlineDir, "trustline_in.txt")
}

// GetSyncCounter retrieves the sync_counter_in value using the datagram to determine the directory.
func GetSyncCounter(dg *types.Datagram) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
	return database.GetUint32FromFile(trustlineDir, "sync_counter.txt")
}

// GetSyncIn retrieves the sync_in value using the datagram to determine the directory.
func GetSyncIn(dg *types.Datagram) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
	return database.GetUint32FromFile(trustlineDir, "sync_in.txt")
}

// GetSyncOut retrieves the sync_out value using the datagram to determine the directory.
func GetSyncOut(dg *types.Datagram) (uint32, error) {
	trustlineDir := database.GetTrustlineDir(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
	return database.GetUint32FromFile(trustlineDir, "sync_out.txt")
}
