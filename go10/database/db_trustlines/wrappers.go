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

// GetTrustlineOutFromDatagram retrieves the outbound trustline using fields from datagram
func SetTrustlineOutFromDatagram(dg *types.Datagram, value uint32) (uint32, error) {
	return SetTrustlineOut(dg.Username, dg.PeerServerAddress, dg.PeerUsername, value)
}

// GetTrustlineInFromDatagram retrieves the inbound trustline using fields from datagram
func SetTrustlineInFromDatagram(dg *types.Datagram, value uint32) (uint32, error) {
	return SetTrustlineIn(dg.Username, dg.PeerServerAddress, dg.PeerUsername, value)
}
