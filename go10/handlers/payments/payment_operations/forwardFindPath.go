package payment_operations

import (
    "encoding/binary"
    "log"

    "ripple/comm"
    "ripple/commands"
    "ripple/handlers"
    "ripple/types"
    "ripple/database/db_pathfinding"
)

// ForwardFindPath forwards the pathfinding request to all connected peers
func ForwardFindPath(datagram *types.Datagram, inOrOut byte) {
    // Retrieve the list of connected peers
    peers, err := db_pathfinding.GetPeers(datagram.Username)
    if err != nil {
        log.Printf("Failed to retrieve peers for user %s: %v", datagram.Username, err)
        return
    }

    amount := binary.BigEndian.Uint32(datagram.Arguments[32:36])

    for _, peer := range peers {
        // Skip if this peer is the one from which the datagram was received
        if peer.Username == datagram.PeerUsername && peer.ServerAddress == datagram.PeerServerAddress {
            continue
        }

        // Check if the trustline is sufficient
        sufficient, err := payments.CheckTrustlineSufficient(datagram.Username, peer.ServerAddress, peer.Username, amount, inOrOut)
        if err != nil {
            log.Printf("Error checking trustline: %v", err)
            continue
        }
        if !sufficient {
            log.Printf("Trustline insufficient for user %s with peer %s at %s", datagram.Username, peer.Username, peer.ServerAddress)
            continue
        }

        err := handlers.PrepareAndSendDatagram(datagram.Command, datagram.Username, peer.ServerAddress, peer.Username, datagram.Arguments[:])
        if err != nil {
            log.Printf("Failed to prepare and send datagram: %v", err)
            return
        }

        log.Printf("Sent pathfinding request to %s at %s", peer.Username, peer.ServerAddress)
    }
}
