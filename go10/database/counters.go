package database

import (
	"ripple/types"
)

// GetCounterIn retrieves the counter_in value using the datagram to determine the directory.
func GetCounterIn(username, peerServerAddress, peerUsername string) (uint32, error) {
	peerDir := GetPeerDir(username, peerServerAddress, peerUsername)
	return GetUint32FromFile(peerDir, "counter_in.txt")
}

// SetCounterIn sets the counter_in value.
func SetCounterIn(username, peerServerAddress, peerUsername string, value uint32) error {
	peerDir := GetPeerDir(username, peerServerAddress, peerUsername)
	return WriteUint32ToFile(peerDir, "counter_in.txt", value)
}

// GetCounterOut retrieves the counter_out value using the datagram to determine the directory.
func GetCounterOut(username, peerServerAddress, peerUsername string) (uint32, error) {
	peerDir := GetPeerDir(username, peerServerAddress, peerUsername)
	return GetUint32FromFile(peerDir, "counter_out.txt")
}

// SetCounterOut sets the counter_out value.
func SetCounterOut(username, peerServerAddress, peerUsername string, value uint32) error {
	peerDir := GetPeerDir(username, peerServerAddress, peerUsername)
	return WriteUint32ToFile(peerDir, "counter_out.txt", value)
}

// wrappers

// GetCounterInFromDatagram retrieves the counter_in value using the incoming datagram
func GetCounterInFromDatagram(dg *types.Datagram) (uint32, error) {
	return GetCounterIn(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
}

// SetCounterInFromDatagram sets the counter_in value from the incoming datagram
func SetCounterInFromDatagram(dg *types.Datagram) error {
	return SetCounterIn(dg.Username, dg.PeerServerAddress, dg.PeerUsername, dg.Counter)
}

// GetCounterOutFromDatagram retrieves the counter_out value using the outgoing datagram
func GetCounterOutFromDatagram(dg *types.Datagram, peerServerAddress string) (uint32, error) {
	return GetCounterOut(dg.PeerUsername, peerServerAddress, dg.Username)
}

// SetCounterOutFromDatagram sets the counter_out value from the outgoing datagram
func SetCounterOutFromDatagram(dg *types.Datagram, peerServerAddress string, value uint32) error {
	return SetCounterOut(dg.PeerUsername, peerServerAddress, dg.Username, value)
}

