package pathfinding

import (
    "log"
)

// PathFindingOut handles the pathfinding output command for a given session
func PathFindingOut(session Session) {
    // Extract the username and identifier from the session's datagram
    username := session.Datagram.Username
    var identifier [32]byte
    copy(identifier[:], session.Datagram.Arguments[:32])
    
    // Attempt to find the account node and path entry
    accountNode := session.PathManager.FindAccount(username)
    
    // Check if the account node exists and find the path entry
    var pathEntry *PathEntry
    if accountNode != nil {
        pathEntry = accountNode.FindPathEntry(identifier)
    }

    // If no path entry exists, terminate the handler as no path has been initiated
    if pathEntry == nil {
        log.Printf("No path entry found for identifier %x and user %s. Handler terminates.", identifier, username)
        return // Early exit if no path entry exists
    }

    // Evaluate the existence of the path entry
    if pathEntry != nil {
        // Path entry exists

        counter := session.Datagram.Counter
        if pathEntry.CounterIn > counter {
    		log.Println("Received counter is not greater than or equal to pathEntry.CounterIn. Potential replay attack.")
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
