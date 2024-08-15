package db_server

import (
	"ripple/database"
)

// GetCounterOut retrieves the counter_out value using the datagram to determine the directory.
func GetCounterOut(dg *Datagram) (uint32, error) {
	peerDir := database.GetPeerDir(dg)
	return database.GetUint32FromFile(peerDir, "counter_out.txt")
}

// GetCounterIn retrieves the counter_in value using the datagram to determine the directory.
func GetCounterIn(dg *Datagram) (uint32, error) {
	peerDir := database.GetPeerDir(dg)
	return database.GetUint32FromFile(peerDir, "counter_in.txt")
}

// SetCounterOut sets the counter_out value.
func SetCounterOut(dg *Datagram, value uint32) error {
	peerDir := database.GetPeerDir(dg)
	return database.WriteUint32ToFile(peerDir, "counter_out.txt", value)
}

// SetCounterIn sets the counter_in value.
func SetCounterIn(dg *Datagram, value uint32) error {
	peerDir := database.GetPeerDir(dg)
	return database.WriteUint32ToFile(peerDir, "counter_in.txt", value)
}
