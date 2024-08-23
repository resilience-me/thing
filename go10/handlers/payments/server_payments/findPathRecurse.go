package server_payments

import (
    "encoding/binary"
    "log"

    "ripple/comm"
    "ripple/commands"
    "ripple/handlers"
    "ripple/pathfinding"
    "ripple/types"
    "ripple/payments"
)

// FindPathRecurse processes a pathfinding recurse command
func FindPathRecurse(session *Session) {
    datagram := session.Datagram

    // Inline extraction of the path identifier and depth from datagram arguments
    pathIdentifier := BytesToString(datagram.Arguments[:32]) // Assuming identifier is in the first 32 bytes
    incomingDepth := binary.BigEndian.Uint32(datagram.Arguments[32:36]) // Assuming depth is in bytes 32-36

    // Find the account using the username from the datagram
    account := pathfinding.PathManager.Find(datagram.Username)
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
    payments_operations.FindPathRecurse(datagram, targetPeer, path.Depth)
}
