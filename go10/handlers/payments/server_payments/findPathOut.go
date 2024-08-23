package main

import (
    "encoding/binary"
    "log"

    "ripple/comm"
    "ripple/commands"
    "ripple/handlers"
    "ripple/pathfinding"
    "ripple/payments"
    "ripple/payments_operations"
    "ripple/types"
    "ripple/database/db_pathfinding"
)

// FindPathOut processes a pathfinding request from the buyer to the seller
func FindPathOut(session *Session) {
    datagram := session.Datagram

    // Inline extraction of the path identifier and amount from datagram arguments
    pathIdentifier := BytesToString(datagram.Arguments[:32]) // Assuming identifier is in the first 32 bytes
    pathAmount := binary.BigEndian.Uint32(datagram.Arguments[32:36]) // Assuming amount is in the next 4 bytes

    // Check if there is sufficient outgoing trustline for the path amount
    sufficient, err := payments.CheckTrustlineSufficient(datagram.Username, datagram.PeerServerAddress, datagram.PeerUsername, pathAmount, types.Outgoing)
    if err != nil {
        log.Printf("Error checking outgoing trustline: %v", err)
        return
    }
    if !sufficient {
        log.Printf("Insufficient outgoing trustline for user %s with peer %s at %s for amount: %d", datagram.Username, datagram.PeerUsername, datagram.PeerServerAddress, pathAmount)
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

        // Since this is the first time seeing this path, send a PathFindingRecurse back to the buyer
        payments_operations.FindPathRecurse(datagram, path.Incoming, 0)
        return
    }

    // If the path is already present, forward the PathFinding request to peers
    log.Printf("Path already exists for identifier %s, forwarding to peers", pathIdentifier)
    payments_operations.FindPath(datagram, types.Outgoing)
}
