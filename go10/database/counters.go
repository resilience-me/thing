package database

import (
	"ripple/types"
)

// GetCounterIn retrieves the counter_in value using the datagram to determine the directory.
func GetCounterIn(dg *types.Datagram) (uint32, error) {
	peerDir := GetPeerDir(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
	return GetUint32FromFile(peerDir, "counter_in.txt")
}

// SetCounterIn sets the counter_in value.
func SetCounterIn(dg *types.Datagram) error {
	peerDir := GetPeerDir(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
	return WriteUint32ToFile(peerDir, "counter_in.txt", dg.Counter)
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

// GetCounterOutFromDatagram retrieves the counter_out value using the outgoing datagram
func GetCounterOutFromDatagram(dg *types.Datagram, peerServerAddress string) (uint32, error) {
	return GetCounterOut(dg.PeerUsername, peerServerAddress, dg.Username)
}

// SetCounterOutFromDatagram sets the counter_out value from the outgoing datagram
func SetCounterOutFromDatagram(dg *types.Datagram, peerServerAddress string, value uint32) error {
	return SetCounterOut(dg.PeerUsername, peerServerAddress, dg.Username, value)
}

// complex function

// GetAndIncrementCounterOut retrieves the current counter_out, increments it, and updates the database.
// It returns the counter value before it was incremented.
func GetAndIncrementCounterOut(username, peerServerAddress, peerUsername string) (uint32, error) {
    // Retrieve the current value of counter_out from the database.
    counterOut, err := GetCounterOut(username, peerServerAddress, peerUsername)
    if err != nil {
        return 0, err  // Return error if unable to fetch the counter.
    }

    // Increment the counter and update it in the database within the same function call.
    if err := SetCounterOut(username, peerServerAddress, peerUsername, counterOut + 1); err != nil {
        return 0, err  // Return error if unable to update the counter.
    }

    // Return the original counter value that was fetched.
    return counterOut, nil
}

// wrapper

func GetAndIncrementCounterOutFromDatagram(dg *types.Datagram, peerServerAddress string) (uint32, error) {
	return GetAndIncrementCounterOut(dg.PeerUsername, peerServerAddress, dg.Username)
}
