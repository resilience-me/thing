package payments

import (
    "fmt"
    "ripple/commands"
)

// CheckTrustlineSufficient checks if the trustline (either incoming or outgoing) is sufficient for the given amount.
func CheckTrustlineSufficient(username, peerServerAddress, peerUsername string, amount uint32, command byte) (bool, error) {
    switch command {
    case commands.ServerPayments_FindPathOut: // Assuming FindPathOut uses incoming trustline
        return CheckTrustlineInSufficient(username, peerServerAddress, peerUsername, amount)
    case commands.ServerPayments_FindPathIn: // Assuming FindPathIn uses outgoing trustline
        return CheckTrustlineOutSufficient(username, peerServerAddress, peerUsername, amount)
    default:
        return false, fmt.Errorf("unsupported command for trustline check: %d", command)
    }
}
