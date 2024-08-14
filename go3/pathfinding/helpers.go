package pathfinding

import (
    "fmt"
)

func (pm *PathManager) cleanupAccounts() {
    now := time.Now()

    for username, account := range pm.Accounts {
        // Only check the LastModified timestamp to decide if the account should be removed
        if now.Sub(account.LastModified) > config.PathFindingTimeout {
            delete(pm.Accounts, username)
        }
    }
}

// initiatePayment sets up or updates payment details for an account, creating the account if necessary.
func (pm *PathManager) initiatePayment(username, identifier string, inOrOut bool) error {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    // Perform account cleanup before processing the new payment
    pm.cleanupAccounts()

    account, exists := pm.Accounts[username]
    if !exists {
        // If the account does not exist, create it and set LastModified to now
        account = &Account{
            Username:      username,
            LastModified:  time.Now(),
            Paths:         make(map[string]*Path),
        }
        pm.Accounts[username] = account
    } else {
        // If account exists, check for an existing payment and handle path removal
        if account.Payment != nil {
            if _, ok := account.Paths[account.Payment.Identifier]; ok {
                delete(account.Paths, account.Payment.Identifier)
            }
        }
        // Update the LastModified only if account already existed and is being modified
        account.LastModified = time.Now()
    }

    // Set or update the payment details
    account.Payment = &Payment{
        Identifier: identifier,
        InOrOut:    inOrOut,
    }

    // Optionally, add or update the related Path entry with a new timestamp
    account.Paths[paymentDetails.Identifier] = &Path{
        Identifier:   identifier,
        Timestamp:    time.Now(),  // This timestamp represents the payment time
        Incoming:     PeerAccount{}, // These would be set according to your logic
        Outgoing:     PeerAccount{},
        CounterIn:    0,
        CounterOut:   make(map[string]int),
    }

    return nil
}

// Wrapper for initiating an outgoing payment
func (pm *PathManager) InitiateOutgoingPayment(username, paymentID string) error {
    return pm.initiatePayment(username, paymentID, true)
}

// Wrapper for initiating an incoming payment
func (pm *PathManager) InitiateIncomingPayment(username, paymentID string) error {
    return pm.initiatePayment(username, paymentID, false)
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


