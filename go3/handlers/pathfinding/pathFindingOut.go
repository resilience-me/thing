package pathfinding

import (
    "log"
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
        pathEntry = accountNode.FindPathEntry(identifier)
    } else {
        // Create a new account node
        accountNode = session.PathManager.AddAccount(username)
        log.Printf("Created new account node for user %s.\n", username)
    }

    // Evaluate the existence of the path entry
    if pathEntry != nil {
        // Path entry exists

        counter := session.Datagram.Counter
        if pathEntry.Counter > counter {
    		log.Println("Received counter is not greater than or equal to pathEntry.Counter. Potential replay attack.")
            return
        }

        // Send the PathFindingRecurse command back to the peer
        responseDatagram := &main.Datagram{
            Command:           main.Pathfinding_PathFindingRecurse,
            Username:          datagram.PeerUsername,             // Send to the peer username
            PeerUsername:      datagram.Username,                  // Original sender as PeerUsername
            PeerServerAddress: config.GetServerAddress(),          // Use config to get the server address
            Arguments:         datagram.Arguments,                 // Include the original Arguments
            Counter:           
        }

        // Log the sending action
        log.Printf("Sending PathFindingRecurse command from %s to %s", datagram.PeerUsername, datagram.Username)


    } else {
        // Path entry does not exist
        return
    }
}
