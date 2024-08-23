package server_payments

import (
    "encoding/binary"
    "log"

    "ripple/comm"
    "ripple/handlers"
    "ripple/pathfinding"
    "ripple/payments"
    "ripple/types"
    "ripple/database/db_pathfinding"
)

// StartFindPath initiates a pathfinding request to all connected peers.
func StartFindPath(username, identifier string, amount uint32, inOrOut byte) {
    peers, err := db_pathfinding.GetPeers(username)
    if err != nil {
        log.Printf("Failed to retrieve peers for user %s: %v", username, err)
        return
    }

    arguments := append([]byte{inOrOut}, []byte(identifier)...)
    arguments = append(arguments, types.Uint32ToBytes(amount)...)

    for _, peer := range peers {
        sufficient, err := payments.CheckTrustlineSufficient(username, peer.ServerAddress, peer.Username, amount, inOrOut)
        if err != nil {
            log.Printf("Error checking trustline: %v", err)
            continue
        }
        if !sufficient {
            log.Printf("Trustline insufficient for user %s with peer %s at %s", username, peer.Username, peer.ServerAddress)
            continue
        }

        newDatagram, err := handlers.PrepareDatagram(commands.ServerPayments_FindPath, username, peer.ServerAddress, peer.Username, arguments)
        if err != nil {
            log.Printf("Failed to prepare datagram: %v", err)
            continue
        }

        if err := comm.SignAndSendDatagram(newDatagram, peer.ServerAddress); err != nil {
            log.Printf("Failed to send pathfinding request to %s at %s: %v", peer.Username, peer.ServerAddress, err)
            return
        }

        log.Printf("Sent pathfinding request to %s at %s", peer.Username, peer.ServerAddress)
    }
}
