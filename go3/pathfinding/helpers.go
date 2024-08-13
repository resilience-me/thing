package pathfinding

import (
    "fmt"
)

// FindOrAdd looks for an existing account and returns it, or adds a new one if it does not exist.
func (pm *PathManager) FindOrAdd(username string) *AccountNode {
    // Attempt to find the existing node
    existingNode := pm.Find(username)
    if existingNode != nil {
        return existingNode
    }

    // Only reach here if no existing node was found; add a new one
    return pm.Add(username)
}

// Shared function to initialize a payment, based on whether it is incoming or outgoing
func (pm *PathManager) initiatePayment(username, paymentID string, inOrOut bool) error {
    // Step 1: Find or add the AccountNode
    accountNode := pm.FindOrAdd(username)
    
    // Step 2: Check if a PathNode for this payment already exists
    pathNode := accountNode.Find(paymentID)
    if pathNode != nil {
        // If a PathNode already exists for this payment, return an error or handle accordingly
        return fmt.Errorf("payment with ID %s is already in progress", paymentID)
    }

    // Step 3: If no PathNode exists, create a new one
    pathNode = accountNode.Add(paymentID, PeerAccount{}, PeerAccount{})

    // Step 4: Initialize the Payment struct and assign it to the AccountNode
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
