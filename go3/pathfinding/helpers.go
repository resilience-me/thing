package pathfinding

import (
    "fmt"
)

// FindOrAdd function finds the account node and resets the timestamp
func (pm *PathManager) FindOrAdd(username string) *AccountNode {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    // Use BaseList.Find to locate the base node directly
    existingNode := pm.BaseList.Find(username)

    if existingNode != nil {
        existingNode.Timestamp = time.Now() // Directly reset the timestamp in the BaseNode
        return existingNode.(*AccountNode)  // Type assert only for returning the correct type
    }

    // Only reach here if no existing node was found; add a new one
    return pm.Add(username)
}

// ResetPayment safely removes the previous payment and clears the Payment field
func (pm *PathManager) ResetPayment(accountNode *AccountNode) {
    // Lock the mutex to ensure thread-safe modification of the AccountNode
    pm.mu.Lock()
    defer pm.mu.Unlock()

    if accountNode.Payment != nil {
        // Remove the previous payment's PathNode
        accountNode.PathList.Remove(accountNode.Payment.Identifier)
        accountNode.Payment = nil // Clear the previous payment
    }
}

// Shared function to initialize a payment, based on whether it is incoming or outgoing
func (pm *PathManager) initiatePayment(username, paymentID string, inOrOut bool) error {
    // Step 1: Find or add the AccountNode
    accountNode := pm.FindOrAdd(username)

    // Step 2: Safely clear the previous payment using the helper function
    pm.ResetPayment(accountNode)

    // Step 3: Check if a PathNode for this payment already exists
    pathNode := accountNode.Find(paymentID)
    if pathNode != nil {
        // If a PathNode already exists for this payment, return an error or handle accordingly
        return fmt.Errorf("payment with ID %s is already in progress", paymentID)
    }

    // Step 4: If no PathNode exists, create a new one
    pathNode = accountNode.Add(paymentID, PeerAccount{}, PeerAccount{})

    // Step 5: Initialize the Payment struct and assign it to the AccountNode
    accountNode.Payment = &Payment{
        Identifier: paymentID,
        InOrOut:    inOrOut, // Set based on whether this is an incoming or outgoing payment
    }

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
