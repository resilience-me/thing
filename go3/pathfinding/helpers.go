package pathfinding

import (
    "time"        // For using time.Now() and time.Duration
    "sync"        // For using sync.Mutex
    "ripple/config" // For using config.PathFindingTimeout
)

// Reinsert updates LastModified and reinserts the account if it was removed.
func (pm *PathManager) reinsert(username string, account *Account) {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    // Update LastModified
    account.LastModified = time.Now()

    // Reinsert the account
    pm.Accounts[username] = account
}

func (pm *PathManager) cleanupAccounts() {
    now := time.Now()

    for username, account := range pm.Accounts {
        if now.After(account.Cleanup) {
            // Logic to handle cleanup of the account
            delete(pm.Accounts, username)
        }
    }
}

// initiatePayment sets up or updates payment details for an account, creating the account if necessary.
func (pm *PathManager) initiatePayment(username, identifier string, inOrOut bool) error {
    // Perform account cleanup before processing the new payment
    pm.cleanupAccounts()
    // Retrieve the account from the PathManager, creating it if necessary
    account := pm.Find(username)
    if account == nil {
        // Account does not exist; create it
        account = pm.Add(username)
    }

    // Remove the existing path if present
    if account.Payment != nil {
        account.Remove(account.Payment.Identifier)
    }

    // Set or update the payment details
    account.Payment = &Payment{
        Identifier: identifier,
        InOrOut:    inOrOut,
    }

    // Add or update the related Path entry with a new timestamp
    account.Add(identifier, PeerAccount{}, PeerAccount{}) // Adjust PeerAccount as needed

    // Update LastModified. Then reinsert, to prevent a minimal race condition risk.
    // The risk is that LastModified could have been very close to timing out when
    // the cleanup was run at the beginning of this function, and another thread
    // could have run the cleanup a few nanoseconds later and deleted the account
    // This is easily mitigated by simply reinserting it, hence reinsert is used.
    pm.reinsert(username, account)

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
