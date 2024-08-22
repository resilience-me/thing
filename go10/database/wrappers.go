// package database

// import "ripple/types"

// // GetAccountDirFromDatagram constructs the account directory path from the datagram
// func GetAccountDirFromDatagram(dg *types.Datagram) string {
//     return GetAccountDir(dg.Username)
// }

// // GetPeerDirFromIncomingDatagram constructs the peer directory path from the incoming datagram and returns it
// func GetPeerDirFromIncomingDatagram(dg *types.Datagram) string {
//     return GetPeerDir(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
// }

// // GetPeerDirFromOutgoingDatagram constructs athe peer directory path from the outgoing datagram and returns it
// func GetPeerDirFromOutgoingDatagram(dg *types.Datagram, peerServerAddress string) string {
//     return GetPeerDir(dg.PeerUsername, peerServerAddress, dg.Username)
// }

// // Wrapper for GetPeerDirFromIncomingDatagram
// func GetPeerDirFromDatagram(dg *types.Datagram) string {
//     return GetPeerDirFromIncomingDatagram(dg)
// }
