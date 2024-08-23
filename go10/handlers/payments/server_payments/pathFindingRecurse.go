func PathFindingRecurse(session *Session) {
    datagram := session.Datagram

  // Inline extraction of the path identifier and depth from datagram arguments
    pathIdentifier := BytesToString(datagram.Arguments[:32]) // Assuming identifier is in the first 32 bytes
    depth := binary.BigEndian.Uint32(datagram.Arguments[32:36]) // Assuming depth is in bytes 32-36

    // Find the account using the username from the datagram
    account := session.pm.Find(datagram.Username)
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

    // Validate the depth
    if depth != path.Depth {
        log.Printf("Depth mismatch for path %s: expected %d, got %d", pathIdentifier, path.Depth, depth)
        return
    }

    // Increment the depth since it matches
    path.Depth++
    log.Printf("Incremented depth for path %s: new depth is %d", pathIdentifier, path.Depth)

    // Proceed with further pathfinding logic or response handling
    // (This could include forwarding the request to peers, sending back a response, etc.)
}
