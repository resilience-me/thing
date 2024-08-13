package pathfinding

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

// InitiatePayment initializes a payment, ensuring the AccountNode and PathNode are correctly set up.
func (pm *PathManager) InitiatePayment(username, paymentID string) error {
    // Step 1: Find or add the AccountNode
    accountNode := pm.FindOrAdd(username)
    
    // Step 2: Check if a PathNode for this payment already exists
    pathNode := accountNode.Find(paymentID)
    if pathNode != nil {
        // If a PathNode already exists for this payment, return an error or handle accordingly
        return fmt.Errorf("payment with ID %s is already in progress", paymentID)
    }

    // Step 3: If no PathNode exists, create a new one
    pathNode = accountNode.Add(paymentID, PeerAccount{}, PeerAccount{}) // Empty PeerAccounts since they aren't parameters

    // Step 4: Initialize the Payment struct and assign it to the AccountNode
    accountNode.Payment = &Payment{
        Identifier: paymentID,
        InOrOut:    true, // Assuming this node is the initiator of the payment
        Deadline:   time.Now().Add(5 * time.Minute), // Assuming a 5-minute deadline for the payment
    }

    accountNode.ActivePaymentPathNode = pathNode // Set the active payment PathNode

    return nil // Indicate success
}
