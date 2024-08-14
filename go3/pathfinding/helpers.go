package pathfinding

import (
    "fmt"
)

// InitiatePayment starts a new payment process by creating or updating a path node.
func (pm *PathManager) InitiatePayment(username, identifier string, peer PeerAccount, inOrOut bool) error {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    account, exists := pm.Accounts[username]
    if !exists {
        return fmt.Errorf("account %s does not exist", username)
    }

    if _, exists := account.Paths[identifier]; exists {
        return fmt.Errorf("path %s already exists for account %s", identifier, username)
    }

    account.Paths[identifier] = &PathNode{
        Identifier:   identifier,
        Incoming:     PeerAccount{},
        Outgoing:     peer,
        LastModified: time.Now(),
    }
    return nil
}

// findOrAdd function finds the account node and resets the timestamp
func (pm *PathManager) findOrAdd(username string) *AccountNode {
    existingNode := pm.Find(username)
    if existingNode != nil {
        return existingNode
    }
    // Only reach here if no existing node was found; add a new one
    return pm.Add(username)
}

// resetPayment removes the previous payment and clears the Payment field
func (pm *PathManager) resetPayment(accountNode *AccountNode) {
    if accountNode.Payment != nil {
        // Remove the previous payment's PathNode
        accountNode.Remove(accountNode.Payment.Identifier)
        accountNode.Payment = nil // Clear the previous payment
    }
}

func (node *AccountNode) newPayment(paymentID string, inOrOut bool) {
    // Initialize the Payment struct and assign it to the AccountNode
    node.Payment = &Payment{
        Identifier: paymentID,
        InOrOut:    inOrOut, // Set based on whether this is an incoming or outgoing payment
    }
    // Create a new path node
    pathNode = node.Add(paymentID, PeerAccount{}, PeerAccount{})
}

// Shared function to initialize a payment, based on whether it is incoming or outgoing
func (pm *PathManager) initiatePayment(username, paymentID string, inOrOut bool) error {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    existingNode := pm.Find(username)
    if existingNode != nil {
        // Check if a PathNode for this payment already exists
        pathNode := accountNode.Find(paymentID)
        if pathNode != nil {
            // If a PathNode already exists for this payment, return an error or handle accordingly
            return fmt.Errorf("payment with ID %s is already in progress", paymentID)
        }
        // clear the previous payment using the helper function
        pm.ResetPayment(accountNode)

        accountNode.Timestamp = time.Now()
        accountNode.newPayment(paymentID, inOrOut)
        
    } else {
        // Find or add the AccountNode
        accountNode := pm.FindOrAdd(username)
        accountNode.newPayment(paymentID, inOrOut)
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
