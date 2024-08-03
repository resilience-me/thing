// package handlers

// import (
//     "net"
//     "os"
//     "path/filepath"
//     "resilience/main"
// )

// // Handles fetching the outbound trustline information
// func GetTrustlineOut(dg main.Datagram, conn *net.UDPConn, addr *net.UDPAddr) {
//     peerDir, err := main.GetPeerDir(dg)
//     if err != nil {
//         return // Optionally handle or log the error
//     }

//     trustlineOutPath := filepath.Join(peerDir, "trustline", "trustline_out.txt")
//     trustlineAmount, err := os.ReadFile(trustlineOutPath)
//     if err != nil {
//         return // Optionally handle or log the error
//     }

//     _, err = conn.WriteToUDP(trustlineAmount, addr)
//     if err != nil {
//         return // Optionally handle or log the error
//     }
// }

// // Handles fetching the inbound trustline information
// func GetTrustlineIn(dg main.Datagram, conn *net.UDPConn, addr *net.UDPAddr) {
//     peerDir, err := main.GetPeerDir(dg)
//     if err != nil {
//         return // Optionally handle or log the error
//     }

//     trustlineInPath := filepath.Join(peerDir, "trustline", "trustline_in.txt")
//     trustlineAmount, err := os.ReadFile(trustlineInPath)
//     if err != nil {
//         return // Optionally handle or log the error
//     }

//     _, err = conn.WriteToUDP(trustlineAmount, addr)
//     if err != nil {
//         return // Optionally handle or log the error
//     }
// }
