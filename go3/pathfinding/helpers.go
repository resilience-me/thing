package pathfinding

import (
    "time"        // For using time.Now() and time.Duration
    "sync"        // For using sync.Mutex
    "ripple/config" // For using config.PathFindingTimeout
)

// Touch checks if an account exists, updates its LastModified timestamp if it does, and returns the account.
func (pm *PathManager) Touch(username string) *Account {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    account, exists := pm.Accounts[username]
    if exists {
        // Update the LastModified timestamp if the account exists
        account.LastModified = time.Now()
    }
    return account
}

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
    // Perform account cleanup before processing the new payment
    pm.cleanupAccounts()

    // Check if the account exists and refresh LastModified if so
    account := pm.Touch(username)

    if account != nil {
        // If account exists, check for an existing payment and handle path removal
        if account.Payment != nil {
            account.Remove(account.Payment.Identifier)
        }
    } else {
        // If the account does not exist, create it
        account = pm.Add(username)
    }

    // Set or update the payment details
    account.Payment = &Payment{
        Identifier: identifier,
        InOrOut:    inOrOut,
    }

    // Add or update the related Path entry with a new timestamp
    account.Add(identifier, PeerAccount{}, PeerAccount{})

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
