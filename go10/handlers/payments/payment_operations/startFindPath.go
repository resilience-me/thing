package server_payments

import (
    "log"

    "ripple/pathfinding"
    "ripple/payments"
    "ripple/types"
    "ripple/database/db_pathfinding"
    "ripple/handlers"
)

// StartFindPath initiates a pathfinding request to all connected peers.
func StartFindPath(username, identifier string, amount uint32, inOrOut byte) {
    // Retrieve the list of connected peers
    peers, err := db_pathfinding.GetPeers(username)
    if err != nil {
        log.Printf("Failed to retrieve peers for user %s: %v", username, err)
        return
    }

    arguments := append([]byte(identifier), types.Uint32ToBytes(amount)...)
    command := GetFindPathCommand(inOrOut)

    for _, peer := range peers {
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

        // Prepare, sign, and send the datagram using the helper function from the handlers package
        err = handlers.PrepareAndSendDatagram(command, username, peer.ServerAddress, peer.Username, arguments)
        if err != nil {
            log.Printf("Failed to prepare and send pathfinding request from %s to peer %s at server %s: %v", username, peer.Username, peer.ServerAddress, err)
            return // Exit early on error
        }

        log.Printf("Sent pathfinding request to %s at %s", peer.Username, peer.ServerAddress)
    }
}
