package payments

import (
    "encoding/binary"
    "log"

    "ripple/comm"
    "ripple/commands"
    "ripple/handlers"
    "ripple/payments"
    "ripple/types"
    "ripple/database/db_pathfinding"
)

func forwardFindPath(datagram *types.Datagram) {
    // Retrieve the list of connected peers
    peers, err := db_pathfinding.GetPeers(datagram.Username)
    if err != nil {
        log.Printf("Failed to retrieve peers for user %s: %v", datagram.Username, err)
        return
    }

    // Extract the path amount from the datagram arguments
    pathAmount := binary.BigEndian.Uint32(datagram.Arguments[32:36])

    for _, peer := range peers {
        // Skip if this peer is the one from which the datagram was received
        if peer.Username == datagram.PeerUsername && peer.ServerAddress == datagram.PeerServerAddress {
            continue
        }

        // Check if the trustline (in or out) is sufficient
        sufficient, err := payments.CheckTrustlineSufficient(datagram.Username, peer.ServerAddress, peer.Username, pathAmount, datagram.Command)
        if err != nil {
            log.Printf("Error checking trustline: %v", err)
            continue
        }
        if !sufficient {
            log.Printf("Trustline insufficient for user %s with peer %s at %s", datagram.Username, peer.Username, peer.ServerAddress)
            continue
        }

        // Create the new datagram for the next pathfinding request
        newDatagram, err := handlers.PrepareDatagram(datagram.Username, peer.ServerAddress, peer.Username)
        if err != nil {
            log.Printf("Failed to prepare pathfinding datagram: %v", err)
            continue
        }

        // Set the command for the outgoing pathfinding request
        newDatagram.Command = datagram.Command

        // Copy the identifier and amount from the original datagram's arguments
        copy(newDatagram.Arguments[:], datagram.Arguments[:]) // Copy the full Arguments field

        // Serialize and sign the datagram
        if err := comm.SignAndSendDatagram(newDatagram, peer.ServerAddress); err != nil {
            log.Printf("Failed to send pathfinding request to %s at %s: %v", peer.Username, peer.ServerAddress, err)
            return // Exit early on error
        }

        log.Printf("Sent pathfinding request to %s at %s", peer.Username, peer.ServerAddress)
    }
}
