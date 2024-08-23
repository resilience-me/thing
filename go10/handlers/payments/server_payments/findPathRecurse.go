func FindPathRecurse(session *Session, pm *pathfinding.PathManager) {
    datagram := session.Datagram

    // Inline extraction of the path identifier and depth from datagram arguments
    pathIdentifier := BytesToString(datagram.Arguments[:32]) // Assuming identifier is in the first 32 bytes
    incomingDepth := binary.BigEndian.Uint32(datagram.Arguments[32:36]) // Assuming depth is in bytes 32-36

    // Find the account using the username from the datagram
    account := pm.Find(datagram.Username)
    if account == nil {
        log.Printf("Account not found for user: %s", datagram.Username)
        return
    }

    // Retrieve the Path object using the identifier
    path := account.Find(pathIdentifier)
    if path == nil {
        log.Printf("Path not found for identifier: %s", pathIdentifier)
        return
    }

    // Validate the depth first
    if incomingDepth != path.Depth {
        log.Printf("Depth mismatch for path %s: expected %d, got %d", pathIdentifier, path.Depth, incomingDepth)
        return
    }

    // Increment the depth since it matches
    path.Depth++
    log.Printf("Incremented depth for path %s: new depth is %d", pathIdentifier, path.Depth)

    // Check if a Payment is already associated with this account and identifier
    if account.Payment != nil && account.Payment.Identifier == pathIdentifier {
        log.Printf("Reached the root for path %s, processing payment", pathIdentifier)
        processPayment(account, path) // Implement this function to handle the payment
        return
    }

    // Check if both incoming and outgoing are set, indicating a path has already been found
    if path.Incoming.Username != "" && path.Outgoing.Username != "" {
        log.Printf("Path already found for path %s, ignoring recurse", pathIdentifier)
        return
    }

    // Determine the direction based on which peer account is populated in the Path
    var targetPeer pathfinding.PeerAccount
    if path.Outgoing.Username != "" {
        // Path is moving forward, pass it back to the incoming peer
        targetPeer = path.Incoming
    } else if path.Incoming.Username != "" {
        // Path is moving backward, pass it to the outgoing peer
        targetPeer = path.Outgoing
    } else {
        log.Printf("Unable to determine direction for path %s, both Incoming and Outgoing are empty", pathIdentifier)
        return
    }

    // Forward the command to the appropriate peer
    forwardPathFindingRecurseCommand(datagram, targetPeer)
}

func forwardPathFindingRecurseCommand(datagram *types.Datagram, targetPeer pathfinding.PeerAccount) {
    // Use the PrepareDatagram helper to create the new datagram with incremented counter
    newDatagram, err := handlers.PrepareDatagram(datagram.Username, targetPeer.ServerAddress, targetPeer.Username)
    if err != nil {
        log.Printf("Failed to prepare datagram: %v", err)
        return
    }

    // Set the command from the original datagram
    newDatagram.Command = datagram.Command

    // Copy the arguments from the original datagram
    newDatagram.Arguments = datagram.Arguments

    // Update the depth in the new datagram arguments based on path.Depth
    binary.BigEndian.PutUint32(newDatagram.Arguments[32:36], path.Depth)

    // Sign and send the datagram with low priority
    err = comm.SignAndSendDatagram(newDatagram, targetPeer.ServerAddress)
    if err != nil {
        log.Printf("Failed to sign and send PathFindingRecurse command to %s at %s: %v", targetPeer.Username, targetPeer.ServerAddress, err)
        return
    }

    log.Printf("Successfully signed and sent PathFindingRecurse command to %s at %s", targetPeer.Username, targetPeer.ServerAddress)
}