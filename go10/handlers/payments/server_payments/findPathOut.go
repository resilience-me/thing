func FindPathOut(session *Session) {
    datagram := session.Datagram
    pm := session.pm // Access PathManager from the session

    // Inline extraction of the path identifier and amount from datagram arguments
    pathIdentifier := BytesToString(datagram.Arguments[:32]) // Assuming identifier is in the first 32 bytes
    pathAmount := binary.BigEndian.Uint32(datagram.Arguments[32:36]) // Assuming amount is in the next 4 bytes

    // Find the account using the username from the datagram
    account := pm.Find(datagram.Username)
    if account == nil {
        log.Printf("Account not found for user: %s", datagram.Username)
        return
    }

    // Retrieve the Path object using the identifier
    path := account.Find(pathIdentifier)
    if path == nil {
        // Path is not found, add the new path using the Add method
        incomingPeer := pathfinding.NewPeerAccount(datagram.PeerUsername, datagram.PeerServerAddress)
        path = account.Add(pathIdentifier, pathAmount, incomingPeer, pathfinding.PeerAccount{})
        log.Printf("Initialized new path for identifier: %s with amount: %d", pathIdentifier, pathAmount)

        // Since this is the first time seeing this path, send a PathFindingRecurse back to the buyer
        findPathOutRecurse(datagram, path)
        return
    }

    // If the path is already present, forward the PathFinding request to peers
    log.Printf("Path already exists for identifier %s, forwarding to peers", pathIdentifier)
    forwardFindPathOut(account, path)
}

func findPathOutRecurse(datagram *types.Datagram, path *pathfinding.Path) {
    // Directly target the incoming peer, which represents the direction back to the buyer
    targetPeer := path.Incoming

    // Prepare the datagram for forwarding
    newDatagram, err := handlers.PrepareDatagram(datagram.Username, targetPeer.ServerAddress, targetPeer.Username)
    if err != nil {
        log.Printf("Failed to prepare datagram: %v", err)
        return
    }

    // Set the command directly to ServerPayments_FindPathRecurse
    newDatagram.Command = commands.ServerPayments_FindPathRecurse

    // Copy the path identifier from path.Identifier to the new datagram's Arguments field
    copy(newDatagram.Arguments[:32], []byte(path.Identifier)) // Assuming path.Identifier is at most 32 bytes

    // No need to update the depth in newDatagram.Arguments[32:36] since it is already zero

    // Sign and send the datagram to the target peer
    if err := comm.SignAndSendDatagram(newDatagram, targetPeer.ServerAddress); err != nil {
        log.Printf("Failed to sign and send FindPathRecurse command to %s at %s: %v", targetPeer.Username, targetPeer.ServerAddress, err)
        return
    }

    log.Printf("Successfully signed and sent FindPathRecurse command to %s at %s", targetPeer.Username, targetPeer.ServerAddress)
}

func forwardFindPathOut(datagram *types.Datagram, path *pathfinding.Path) {
    // Retrieve the list of connected peers
    peers, err := db_pathfinding.GetPeers(datagram.Username)
    if err != nil {
        log.Printf("Failed to retrieve peers for user %s: %v", datagram.Username, err)
        return
    }

    // Extract the path amount from the datagram arguments
    pathAmount := binary.BigEndian.Uint32(datagram.Arguments[32:36])

    for _, peer := range peers {
        // Skip if this peer is already part of the path as the incoming peer
        if peer.Username == path.Incoming.Username && peer.ServerAddress == path.Incoming.ServerAddress {
            continue
        }

        // Check if the incoming trustline is sufficient
        trustlineIn, err := db_trustlines.GetTrustlineIn(datagram.Username, peer.ServerAddress, peer.Username)
        if err != nil {
            log.Printf("Failed to retrieve incoming trustline for user %s with peer %s at %s: %v", datagram.Username, peer.Username, peer.ServerAddress, err)
            continue
        }

        if trustlineIn < pathAmount {
            log.Printf("Insufficient incoming trustline for user %s with peer %s at %s. Required: %d, Available: %d", datagram.Username, peer.Username, peer.ServerAddress, pathAmount, trustlineIn)
            continue
        }

        // Create the new datagram for the next pathfinding request
        newDatagram, err := handlers.PrepareDatagram(datagram.Username, peer.ServerAddress, peer.Username)
        if err != nil {
            log.Printf("Failed to prepare pathfinding datagram: %v", err)
            continue
        }

        // Set the command for the outgoing pathfinding request
        newDatagram.Command = commands.ServerPayments_FindPathOut

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