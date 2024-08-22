package database

import "ripple/types"

// GetAccountDirFromDatagram constructs the account directory path from the datagram
func GetAccountDirFromDatagram(dg *types.Datagram) string {
    return GetAccountDir(dg.Username)
}

// GetPeerDirFromDatagram constructs the peer directory path from the datagram and returns it
func GetPeerDirFromDatagram(dg *types.Datagram) string {
    return GetPeerDir(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
}

// GetPeerDirFromOutgoingDatagram constructs the peer directory path from the datagram and returns it
func GetPeerDirFromOutgoingDatagram(dg *types.Datagram, peerServerAddress string) string {
    return GetPeerDir(dg.PeerUsername, peerServerAddress, dg.Username)
}
