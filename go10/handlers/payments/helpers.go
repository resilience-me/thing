package payments

import (
    "fmt"
    "ripple/database/db_trustlines"
)

// CheckTrustlineInSufficient checks if the incoming trustline is sufficient for the given amount.
// It returns true if the trustline is sufficient, false otherwise, along with any error encountered.
func CheckTrustlineInSufficient(username, peerServerAddress, peerUsername string, amount uint32) (bool, error) {
    // Retrieve the incoming trustline
    trustlineIn, err := db_trustlines.GetTrustlineIn(username, peerServerAddress, peerUsername)
    if err != nil {
        return false, fmt.Errorf("failed to retrieve incoming trustline for user %s with peer %s at %s: %v", username, peerUsername, peerServerAddress, err)
    }

    // Check if the trustline is sufficient
    if trustlineIn < amount {
        return false, nil
    }

    return true, nil
}
