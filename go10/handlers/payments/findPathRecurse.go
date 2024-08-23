package payments

import (
    "log"
    "ripple/types"
    "ripple/handlers"
    "ripple/comm"
    "ripple/commands"
    "ripple/pathfinding"
)

// findPathRecurse sends a PathFindingRecurse command to the specified peer.
func findPathRecurse(datagram *types.Datagram, peer pathfinding.PeerAccount) {    
    // Prepare the datagram for forwarding with the command and arguments
    newDatagram, err := handlers.PrepareDatagram(commands.ServerPayments_FindPathRecurse, datagram.Username, peer.ServerAddress, peer.Username, datagram.Arguments[:])
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
