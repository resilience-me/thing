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

    // Send pathfinding requests to all peers, depending on existing counters
    for _, peer := range peers {
        if _, exists := pathNode.CounterOut[peer.Username]; !exists {
            // Send a new pathfinding request if no counter exists for this peer
            err := SendNewPathFindingOut(peer, datagram.Arguments)
            if err != nil {
                log.Printf("Failed to send new pathfinding request to %s: %v", peer.Username, err)
            }
        } else {
            // Update or handle existing pathfinding state
            err := SendPathFindingOut(peer, datagram.Arguments)
            if err != nil {
                log.Printf("Failed to update pathfinding request to %s: %v", peer.Username, err)
            }
        }
    }
}

// SendNewPathFindingOut simulates sending a new pathfinding request
func SendNewPathFindingOut(peer db_pathfinding.PeerAccount, args []byte) error {
    // Simulated function for sending out pathfinding requests
    log.Printf("Sending new pathfinding request to %s at %s", peer.Username, peer.ServerAddress)
    return nil
}

// SendPathFindingOut simulates sending an updated pathfinding request
func SendPathFindingOut(peer db_pathfinding.PeerAccount, args []byte) error {
    // Simulated function for sending out pathfinding requests
    log.Printf("Updating pathfinding request to %s at %s", peer.Username, peer.ServerAddress)
    return nil
}
