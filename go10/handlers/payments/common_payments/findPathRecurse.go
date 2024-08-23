package common_payments

import (
    "log"
    "ripple/types"
    "ripple/handlers"
    "ripple/comm"
    "ripple/commands"
    "ripple/pathfinding"
)

// FindPathRecurse sends a PathFindingRecurse command to the specified peer using the depth and identifier from the datagram.
func FindPathRecurse(datagram *types.Datagram, peer pathfinding.PeerAccount, depth uint32) {
    // Create the arguments slice by appending the depth to the identifier from the datagram
    arguments := append(datagram.Arguments[:32], types.Uint32ToBytes(depth)...)

    // Prepare the datagram for forwarding with the command and updated arguments
    newDatagram, err := handlers.PrepareDatagram(commands.ServerPayments_FindPathRecurse, datagram.Username, peer.ServerAddress, peer.Username, arguments)
    if err != nil {
        log.Printf("Failed to prepare datagram: %v", err)
        return
    }

    // Sign and send the datagram to the target peer
    if err := comm.SignAndSendDatagram(newDatagram, peer.ServerAddress); err != nil {
        log.Printf("Failed to sign and send FindPathRecurse command to %s at %s: %v", peer.Username, peer.ServerAddress, err)
        return
    }

    log.Printf("Successfully signed and sent FindPathRecurse command to %s at %s", peer.Username, peer.ServerAddress)
}
