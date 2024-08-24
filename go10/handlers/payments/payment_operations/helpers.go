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
    creditline, err := db_trustlines.GetCreditline(username, peerServerAddress, peerUsername, inOrOut)
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

// CheckTrustlineAndSendFindPathDatagram checks the trustline and sends the datagram if sufficient.
func CheckTrustlineAndSendFindPathDatagram(command byte, username, peerServerAddress, peerUsername string, amount uint32, inOrOut byte, arguments []byte) error {
    // Check if the trustline is sufficient
    sufficient, err := CheckTrustlineSufficient(username, peerServerAddress, peerUsername, amount, inOrOut)
    if err != nil {
        return fmt.Errorf("error checking trustline: %v", err)
    }
    if !sufficient {
        return fmt.Errorf("trustline insufficient for user %s with peer %s at %s", username, peerUsername, peerServerAddress)
    }

    // Prepare, sign, and send the datagram
    if err := PrepareAndSendDatagram(command, username, peerServerAddress, peerUsername, arguments); err != nil {
        return fmt.Errorf("failed to prepare and send pathfinding request from %s to peer %s at server %s: %v", username, peerUsername, peerServerAddress, err)
    }

    return nil
}
