package payments

import (
    "fmt"
    "ripple/database/db_trustlines"
)

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
