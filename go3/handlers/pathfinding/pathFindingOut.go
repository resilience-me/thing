package pathfinding

import (
    "fmt"
)

// PathFindingOut handles the pathfinding output command for a given session
func PathFindingOut(session Session) {
    // Extract the username from the Datagram
    username := session.Datagram.Username

    // Extract the identifier from the Arguments (assuming it is stored in the first 32 bytes)
    var identifier [32]byte
    copy(identifier[:], session.Datagram.Arguments[:32]) // Adjust based on the actual position

    // Find the account node by username
    accountNode := session.PathManager.FindAccount(username)

    // Check if the accountNode is not nil and then look for the path entry
    var pathEntry *PathEntry
    if accountNode != nil {
        pathEntry = accountNode.FindIdentifier(identifier)
    } else {
        session.PathManager.AddAccount(username)
    }

    // Evaluate the existence of the path entry
    if pathEntry != nil {

    } else {

    }
}
