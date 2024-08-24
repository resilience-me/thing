package server_payments

import (
    "log"

    "ripple/types"
    "ripple/pathfinding"
    "ripple/handlers/payments"
    "ripple/handlers/payments/payment_operations"
)

// PathRecurse processes a pathfinding recurse command
func PathRecurse(session types.Session) {
    datagram := session.Datagram

    // Inline extraction of the path identifier and depth from datagram arguments
    pathIdentifier := types.BytesToString(datagram.Arguments[:32]) // Assuming identifier is in the first 32 bytes
    incomingDepth := types.BytesToUint32(datagram.Arguments[32:36]) // Assuming depth is in bytes 32-36

    // Find the account using the username from the datagram
    account := pathfinding.GetPathManager().Find(datagram.Username)
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
        log.Printf("Reached the root for path %s, sending out new FindPath requests", pathIdentifier)
        // Use the InOrOut field from the Payment object to determine the direction
        payment_operations.StartFindPath(datagram.Username, pathIdentifier, path.Amount, account.Payment.InOrOut)
        return
    }

    // Check if both incoming and outgoing are set, indicating a path has already been found
    if payments.CheckPathFound(path) {
        log.Printf("Path already found for path %s, ignoring recurse", pathIdentifier)
        return
    }
    // Determine the direction based on which peer account is populated in the Path
    targetPeer, err := payments.GetRecursePeer(path)
    if err != nil {
        log.Printf("Error determining target peer: %v", err)
        return
    }

    // Forward the command to the appropriate peer
    payment_operations.PathRecurse(datagram, targetPeer, path.Depth)
}
