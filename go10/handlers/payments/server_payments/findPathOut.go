package server_payments

import (
    "encoding/binary"
    "log"

    "ripple/comm"
    "ripple/handlers"
    "ripple/pathfinding"
    "ripple/payments"
    "ripple/types"
    "ripple/payment_operations"
)

func FindPath(session *Session) {
    datagram := session.Datagram

    // Extract direction (inOrOut) from the first byte of arguments
    inOrOut := datagram.Arguments[0]
    pathIdentifier := BytesToString(datagram.Arguments[1:33]) // Assuming identifier is in bytes 1-32
    pathAmount := binary.BigEndian.Uint32(datagram.Arguments[33:37]) // Assuming amount is in bytes 33-36

    // Check if the trustline is sufficient based on the direction
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
    account := pathfinding.PathManager.Find(datagram.Username)
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

        // Since this is the first time seeing this path, send a PathFindingRecurse back to the origin
        payment_operations.FindPathRecurse(datagram, path.Incoming, 0)
        return
    }

    // If the path is already present, forward the PathFinding request to peers
    log.Printf("Path already exists for identifier %s, forwarding to peers", pathIdentifier)
    payments_operations.ForwardFindPath(datagram, inOrOut)
}
