package db_trustlines

import "ripple/types"

// GetTrustlineOutFromDatagram retrieves the outbound trustline using fields from datagram
func GetTrustlineOutFromDatagram(dg *types.Datagram) (uint32, error) {
	return GetTrustlineOut(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
}

// GetTrustlineInFromDatagram retrieves the inbound trustline using fields from datagram
func GetTrustlineInFromDatagram(dg *types.Datagram) (uint32, error) {
	return GetTrustlineIn(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
}

// SetTrustlineOutFromDatagram sets the outbound trustline amount using fields from datagram
func SetTrustlineOutFromDatagram(dg *types.Datagram, value uint32) (uint32, error) {
	return SetTrustlineOut(dg.Username, dg.PeerServerAddress, dg.PeerUsername, value)
}

// SetTrustlineInFromDatagram sets the inbound trustline amount using fields from datagram
func SetTrustlineInFromDatagram(dg *types.Datagram, value uint32) (uint32, error) {
	return SetTrustlineIn(dg.Username, dg.PeerServerAddress, dg.PeerUsername, value)
}

// GetTrustline retrieves the trustline (either incoming or outgoing) based on the inOrOut parameter.
func GetTrustline(username, peerServerAddress, peerUsername string, inOrOut byte) (uint32, error) {
    if inOrOut == 0 { // Assume 0 means incoming trustline
        return GetTrustlineIn(username, peerServerAddress, peerUsername)
    } else { // Assume 1 means outgoing trustline
        return GetTrustlineOut(username, peerServerAddress, peerUsername)
    }
}
