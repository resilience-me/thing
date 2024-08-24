package server_payments

import (
    "log"
    "ripple/handlers"
    "ripple/payments"
    "ripple/types"
    "ripple/database/db_pathfinding"
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
        // Use the new helper function to check the trustline and send the datagram
        if err := handlers.CheckTrustlineAndSendFindPathDatagram(command, username, peer.ServerAddress, peer.Username, amount, inOrOut, arguments); err != nil {
            log.Printf("Error processing datagram: %v", err)
            continue
        }

        log.Printf("Sent pathfinding request to %s at %s", peer.Username, peer.ServerAddress)
    }
}
