package payments

import (
    "fmt"
    "ripple/database/db_trustlines"
    "ripple/commands"
    "ripple/types"
    "ripple/pathfinding"
)

// GetRecursePeer determines the target peer based on the populated fields in the Path.
func GetRecursePeer(path *pathfinding.Path) (pathfinding.PeerAccount, error) {
    if path.Outgoing.Username != "" {
        // Path is moving forward, pass it back to the incoming peer
        return path.Incoming, nil
    } else if path.Incoming.Username != "" {
        // Path is moving backward, pass it to the outgoing peer
        return path.Outgoing, nil
    }

    return pathfinding.PeerAccount{}, fmt.Errorf("Unable to determine direction for path, both Incoming and Outgoing are empty")
}

// CheckPathFound checks if both incoming and outgoing peers are set in the path, indicating a complete path.
func CheckPathFound(path *pathfinding.Path) bool {
    return path.Incoming.Username != "" && path.Outgoing.Username != ""
}

// DetermineCommand returns the appropriate command based on the inOrOut parameter.
func GetFindPathCommand(inOrOut byte) byte {
    if inOrOut == types.Incoming {
        return commands.ServerPayments_FindPathIn
    }
    return commands.ServerPayments_FindPathOut
}
