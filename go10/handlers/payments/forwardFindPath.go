package payments

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
func forwardFindPath(datagram *types.Datagram, inOrOut byte) {

    username := datagram.Username

    // Retrieve the list of connected peers
    peers, err := db_pathfinding.GetPeers(username)
    if err != nil {
        log.Printf("Failed to retrieve peers for user %s: %v", username, err)
        return
    }

    // Extract datagram fields outside of the loop for efficiency
    command := datagram.Command
    arguments := datagram.Arguments[:]
    
    peerUsername := datagram.PeerUsername
    peerServerAddress := datagram.PeerServerAddress

    // Extract the path amount from the datagram arguments
    amount := binary.BigEndian.Uint32(arguments[32:36])

    for _, peer := range peers {
        // Skip if this peer is the one from which the datagram was received
        if peer.Username == peerUsername && peer.ServerAddress == peerServerAddress {
            continue
        }

        // Check if the trustline is sufficient
        sufficient, err := payments.CheckTrustlineSufficient(username, peer.ServerAddress, peer.Username, amount, inOrOut)
        if err != nil {
            log.Printf("Error checking trustline: %v", err)
            continue
        }
        if !sufficient {
            log.Printf("Trustline insufficient for user %s with peer %s at %s", username, peer.Username, peer.ServerAddress)
            continue
        }

        // Use PrepareDatagram to create the new datagram with command and arguments
        newDatagram, err := handlers.PrepareDatagram(command, username, peer.ServerAddress, peer.Username, arguments)
        if err != nil {
            log.Printf("Failed to prepare datagram: %v", err)
            continue
        }

        // Serialize and sign the datagram
        if err := comm.SignAndSendDatagram(newDatagram, peer.ServerAddress); err != nil {
            log.Printf("Failed to send pathfinding request to %s at %s: %v", peer.Username, peer.ServerAddress, err)
            return // Exit early on error
        }

        log.Printf("Sent pathfinding request to %s at %s", peer.Username, peer.ServerAddress)
    }
}
