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

// CheckTrustlineSufficient checks if the trustline (either incoming or outgoing) is sufficient for the given amount.
func CheckTrustlineSufficient(username, peerServerAddress, peerUsername string, amount uint32, inOrOut byte) (bool, error) {
    // Get the relevant trustline
    trustline, err := db_trustlines.GetTrustline(username, peerServerAddress, peerUsername, inOrOut)
    if err != nil {
        return false, fmt.Errorf("failed to retrieve trustline: %v", err)
    }

    // Get the relevant credit line
    creditline, err := GetCreditline(username, peerServerAddress, peerUsername, inOrOut)
    if err != nil {
        return false, fmt.Errorf("failed to retrieve creditline: %v", err)
    }

    // Calculate the available trustline after accounting for the credit line
    available := trustline - creditline

    // Check if the available trustline is sufficient
    if available < amount {
        return false, nil
    }

    return true, nil
}
