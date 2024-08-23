package server_payments

import (
    "encoding/binary"
    "log"

    "ripple/comm"
    "ripple/commands"
    "ripple/handlers"
    "ripple/pathfinding"
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

    var command byte
    if inOrOut == types.Incoming {
        command = commands.ServerPayments_FindPathIn
    } else {
        command = commands.ServerPayments_FindPathOut
    }

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

        // Create a new datagram for each peer
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
