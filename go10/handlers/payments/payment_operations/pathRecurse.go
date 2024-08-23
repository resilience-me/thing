package payment_operations

import (
    "log"
    "ripple/types"
    "ripple/handlers"
    "ripple/commands"
    "ripple/pathfinding"
)

// PathRecurse sends a PathFindingRecurse command to the specified peer using the depth and identifier from the datagram.
func PathRecurse(datagram *types.Datagram, peer pathfinding.PeerAccount, depth uint32) {
    // Create the arguments slice by appending the depth to the identifier from the datagram
    arguments := append(datagram.Arguments[:32], types.Uint32ToBytes(depth)...)

    // Prepare, sign, and send the datagram using the helper function from the handlers package
    if err := handlers.PrepareAndSendDatagram(commands.ServerPayments_PathRecurse, datagram.Username, peer.ServerAddress, peer.Username, arguments); err != nil {
        log.Printf("Failed to prepare and send PathRecurse command from %s to peer %s at server %s: %v", datagram.Username, peer.Username, peer.ServerAddress, err)
        return
    }

    log.Printf("Successfully signed and sent PathRecurse command from %s to peer %s at server %s", datagram.Username, peer.Username, peer.ServerAddress)
}
