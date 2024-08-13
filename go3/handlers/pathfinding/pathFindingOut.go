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

    // Check if the accountNode exists
    var pathEntry *PathEntry
    if accountNode != nil {
        // Account exists, search for the path entry
        pathEntry = accountNode.FindIdentifier(identifier)
    } else {
        // Create a new account node
        accountNode = session.PathManager.AddAccount(username)
        fmt.Printf("Created new account node for user %s.\n", username)
    }

    // Evaluate the existence of the path entry
    if pathEntry != nil {
        // Path entry exists
    } else {
        // Path entry does not exist
    }
}
