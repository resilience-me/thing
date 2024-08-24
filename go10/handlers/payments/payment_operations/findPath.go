package payment_operations

import (
    "encoding/binary"
    "log"
    "ripple/pathfinding"
    "ripple/handlers/payments"
    "ripple/types"
)

// FindPath handles the common logic for processing FindPath requests.
func FindPath(datagram *types.Datagram, inOrOut byte) {
    // Extract the path identifier and amount from datagram arguments
    pathIdentifier := BytesToString(datagram.Arguments[:32])
    pathAmount := binary.BigEndian.Uint32(datagram.Arguments[32:36])

    // Check if the trustline (incoming or outgoing) is sufficient for the path amount
    sufficient, err := payments.CheckTrustlineSufficient(datagram.Username, datagram.PeerServerAddress, datagram.PeerUsername, pathAmount, inOrOut)
    if err != nil {
        log.Printf("Error checking trustline: %v", err)
        return
    }
    if !sufficient {
        log.Printf("Insufficient trustline for user %s with peer %s at %s for amount: %d", datagram.Username, datagram.PeerUsername, datagram.PeerServerAddress, pathAmount)
        return
    }

    // Find the account using the username from the datagram
    account := pathfinding.GetPathManager().Find(datagram.Username)
    if account == nil {
        log.Printf("Account not found for user: %s", datagram.Username)
        return
    }

    // Retrieve the Path object using the identifier
    path := account.Find(pathIdentifier)
    if path == nil {
        // Path is not found, add the new path using the Add method
        newPeer := pathfinding.NewPeerAccount(datagram.PeerUsername, datagram.PeerServerAddress)
        if inOrOut == types.Outgoing {
            path = account.Add(pathIdentifier, pathAmount, newPeer, pathfinding.PeerAccount{})
        } else {
            path = account.Add(pathIdentifier, pathAmount, pathfinding.PeerAccount{}, newPeer)
        }
        log.Printf("Initialized new path for identifier: %s with amount: %d", pathIdentifier, pathAmount)

        // Send a PathFindingRecurse back to the appropriate peer
        PathRecurse(datagram, newPeer, 0)
        return
    }

    // If the path is already present, forward the PathFinding request to peers
    log.Printf("Path already exists for identifier %s, forwarding to peers", pathIdentifier)
    ForwardFindPath(datagram, inOrOut)
}
