package payments

import (
    "fmt"
)

// CheckTrustlineSufficient checks if the trustline (either incoming or outgoing) is sufficient for the given amount.
func CheckTrustlineSufficient(username, peerServerAddress, peerUsername string, amount uint32, inOrOut byte) (bool, error) {
    if inOrOut == 0 { // Assume 0 means incoming trustline
        return CheckTrustlineInSufficient(username, peerServerAddress, peerUsername, amount)
    } else { // Assume 1 means outgoing trustline
        return CheckTrustlineOutSufficient(username, peerServerAddress, peerUsername, amount)
    }
}
