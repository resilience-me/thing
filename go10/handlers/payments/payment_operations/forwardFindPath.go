package payment_operations

import (
    "encoding/binary"
    "log"
    "ripple/handlers"
    "ripple/payments"
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

        // Use the new CheckAndSendDatagram helper function to handle trustline checking and datagram sending
        if err := handlers.CheckAndSendDatagram(datagram.Command, datagram.Username, peer.ServerAddress, peer.Username, amount, inOrOut, datagram.Arguments[:]); err != nil {
            log.Printf("Failed to process pathfinding request from %s to peer %s at server %s: %v", datagram.Username, peer.Username, peer.ServerAddress, err)
            continue
        }

        log.Printf("Successfully sent pathfinding request from %s to peer %s at server %s", datagram.Username, peer.Username, peer.ServerAddress)
    }
}
