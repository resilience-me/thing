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
        // Path entry does not exist, create a new one
        outgoing := PeerAccount{
            Username:      session.Datagram.PeerUsername,       // Get the peer username from the Datagram
            ServerAddress: session.Datagram.PeerServerAddress,   // Get the peer server address from the Datagram
        }
        incoming := PeerAccount{} // Zero value for Incoming since this is an outgoing request

        // Use the AddPathEntry method to add the new path entry
        accountNode.AddPathEntry(identifier, incoming, outgoing)

        fmt.Printf("Created new path entry for account %s with identifier %x.\n", username, identifier)
    }
}
