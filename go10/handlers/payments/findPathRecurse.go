// findPathRecurse sends a PathFindingRecurse command to the specified peer.
func findPathRecurse(datagram *types.Datagram, targetPeer pathfinding.PeerAccount) {
    // Prepare the datagram for forwarding
    newDatagram, err := handlers.PrepareDatagram(datagram.Username, targetPeer.ServerAddress, targetPeer.Username)
    if err != nil {
        log.Printf("Failed to prepare datagram: %v", err)
        return
    }

    // Set the command to FindPathRecurse
    newDatagram.Command = commands.ServerPayments_FindPathRecurse

    // Copy the path identifier from the original datagram's Arguments field
    copy(newDatagram.Arguments[:32], datagram.Arguments[:32]) // Copy the path identifier (assuming it's in the first 32 bytes)

    // No need to update the depth in newDatagram.Arguments[32:36] since it is already zero

    // Sign and send the datagram to the target peer
    if err := comm.SignAndSendDatagram(newDatagram, targetPeer.ServerAddress); err != nil {
        log.Printf("Failed to sign and send FindPathRecurse command to %s at %s: %v", targetPeer.Username, targetPeer.ServerAddress, err)
        return
    }

    log.Printf("Successfully signed and sent FindPathRecurse command to %s at %s", targetPeer.Username, targetPeer.ServerAddress)
}
