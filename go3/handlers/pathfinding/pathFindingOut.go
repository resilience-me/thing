package pathfinding

import (
    "log"
    "ripple/main"
    "ripple/database/db_pathfinding"
)

// PathFindingOut handles the pathfinding output command for a given session.
func PathFindingOut(session main.Session) {
    datagram := session.Datagram
    username := datagram.Username
    var identifier [32]byte
    copy(identifier[:], datagram.Arguments[:32])
    
    accountNode := session.PathManager.Find(username)
    if accountNode == nil {
        log.Printf("No account node found for user %s. Terminating handler.", username)
        return
    }
    
    pathNode := accountNode.Find(identifier)
    if pathNode == nil {
        log.Printf("No path entry found for identifier %x and user %s. Handler terminates.", identifier, username)
        return
    }

    // Assuming the pathNode.CounterIn needs to be checked against datagram.Counter
    if pathNode.CounterIn >= datagram.Counter {
        log.Println("Received counter is not valid or a replay attack may be happening.")
        return
    }

    // Retrieve all peers associated with this user
    peers, err := db_pathfinding.GetPeers(username)
    if err != nil {
        log.Printf("Error retrieving peers for user %s: %v", username, err)
        return
    }

    // Prepare common datagram fields
    newDatagram := &main.Datagram{
        PeerUsername:      username,  // Sender's username
        PeerServerAddress: config.GetServerAddress(), // Sender's server address
        Arguments:         session.Datagram.Arguments, // Arguments originally received
    }

    // Send pathfinding requests to all peers, depending on existing counters
    for _, peer := range peers {

        // Prepare common datagram fields
        newDatagram := &main.Datagram{
            Username:          peer.Username,
            PeerUsername:      username,  // Sender's username
            PeerServerAddress: config.GetServerAddress(), // Sender's server address
            Arguments:         session.Datagram.Arguments, // Arguments originally received
        }

        if _, exists := pathNode.CounterOut[peer.Username]; !exists {
            newDatagram.Command = main.Pathfinding_NewPathFindingOut
        } else {
            newDatagram.Command = main.Pathfinding_PathFindingOut
        }

        // Sign and send the datagram to the peer
        if err := handlers.SignAndSendDatagram(session, commonDatagram); err != nil {
            fmt.Printf("Failed to send pathfinding datagram to %s: %v\n", peer.Username, err)
            continue // Continue with other peers even if one fails
        }
    }
}
