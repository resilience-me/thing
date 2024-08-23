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

    // Retrieve the outgoing credit line
    creditLineOut, err := db_trustlines.GetCreditLineOut(username, peerServerAddress, peerUsername)
    if err != nil {
        return false, fmt.Errorf("failed to retrieve outgoing credit line for user %s with peer %s at %s: %v", username, peerUsername, peerServerAddress, err)
    }

    // Adjust the trustline by subtracting the outgoing credit line
    availableAmount := trustlineIn - creditLineOut

    // Check if the adjusted trustline is sufficient
    if availableAmount < amount {
        return false, nil
    }

    return true, nil
}

// CheckTrustlineOutSufficient checks if the outgoing trustline is sufficient for the given amount.
// It returns true if the trustline is sufficient, false otherwise, along with any error encountered.
func CheckTrustlineOutSufficient(username, peerServerAddress, peerUsername string, amount uint32) (bool, error) {
    // Retrieve the incoming trustline
    trustlineOut, err := db_trustlines.GetTrustlineOut(username, peerServerAddress, peerUsername)
    if err != nil {
        return false, fmt.Errorf("failed to retrieve incoming trustline for user %s with peer %s at %s: %v", username, peerUsername, peerServerAddress, err)
    }

    // Retrieve the outgoing credit line
    creditLineIn, err := db_trustlines.GetCreditLineIn(username, peerServerAddress, peerUsername)
    if err != nil {
        return false, fmt.Errorf("failed to retrieve outgoing credit line for user %s with peer %s at %s: %v", username, peerUsername, peerServerAddress, err)
    }

    // Adjust the trustline by subtracting the outgoing credit line
    availableAmount := trustlineOut - creditLineIn

    // Check if the adjusted trustline is sufficient
    if availableAmount < amount {
        return false, nil
    }

    return true, nil
}
