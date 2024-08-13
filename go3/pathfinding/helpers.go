package pathfinding

import (
    "fmt"
)

// FindOrAdd function finds the account node and resets the timestamp
func (pm *PathManager) FindOrAdd(username string) *AccountNode {
    existingNode := pm.Find(username)
    if existingNode != nil {
        return existingNode
    }
    // Only reach here if no existing node was found; add a new one
    return pm.Add(username)
}

// ResetPayment removes the previous payment and clears the Payment field
func (pm *PathManager) ResetPayment(accountNode *AccountNode) {
    if accountNode.Payment != nil {
        // Remove the previous payment's PathNode
        accountNode.PathList.Remove(accountNode.Payment.Identifier)
        accountNode.Payment = nil // Clear the previous payment
    }
}

// Shared function to initialize a payment, based on whether it is incoming or outgoing
func (pm *PathManager) initiatePayment(username, paymentID string, inOrOut bool) error {

    pm.mu.Lock()
    defer pm.mu.Unlock()

    // Find or add the AccountNode
    accountNode := pm.FindOrAdd(username)

    // Check if a PathNode for this payment already exists
    pathNode := accountNode.Find(paymentID)
    if pathNode != nil {
        // If a PathNode already exists for this payment, return an error or handle accordingly
        return fmt.Errorf("payment with ID %s is already in progress", paymentID)
    }

    // Safely clear the previous payment using the helper function
    pm.ResetPayment(accountNode)

    // Initialize the Payment struct and assign it to the AccountNode
    accountNode.Payment = &Payment{
        Identifier: paymentID,
        InOrOut:    inOrOut, // Set based on whether this is an incoming or outgoing payment
    }

    // Create a new path node
    pathNode = accountNode.Add(paymentID, PeerAccount{}, PeerAccount{})

    accountNode.Timestamp = time.Now()

    return nil // Indicate success
}

// Wrapper for initiating an outgoing payment
func (pm *PathManager) InitiateOutgoingPayment(username, paymentID string) error {
    return pm.initiatePayment(username, paymentID, true)
}

// Wrapper for initiating an incoming payment
func (pm *PathManager) InitiateIncomingPayment(username, paymentID string) error {
    return pm.initiatePayment(username, paymentID, false)
}
