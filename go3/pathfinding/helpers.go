package pathfinding

import (
    "sync"        // For using sync.Mutex
    "time"        // For using time.Now()
    "ripple/config" // For using config.PathFindingTimeout
)

func (pm *PathManager) cleanupAccounts() {
    now := time.Now()
    for username, account := range pm.Accounts {
        if now.After(account.Cleanup) {
            delete(pm.Accounts, username)
        }
    }
}

// cleanupPaths removes expired paths within the Account.
func (account *Account) cleanupPaths() {
    now := time.Now()
    for pathID, path := range account.Paths {
        if now.After(path.Timeout) {
            delete(account.Paths, pathID)  // Remove expired paths
        }
    }
}

func (pm *PathManager) cleanupCacheAndFetchAccount(username string) *Account {
    // Cleanup all accounts first
    pm.cleanupAccounts()

    account := pm.Find(username)

    if account == nil {
        return pm.Add(username)
    }
    account.cleanupPaths()

    return account
}

// Reinsert updates LastModified and reinserts the account if it was removed.
func (pm *PathManager) reinsert(username string, account *Account) {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    // Update Cleanup field
    account.Cleanup = time.Now().Add(config.PathFindingTimeout)

    // Reinsert the account
    pm.Accounts[username] = account
}

// initiatePayment sets up or updates payment details for an account, creating the account if necessary.
func (pm *PathManager) initiatePayment(username string, payment *Payment) {
    // Fetch or create the account, with any necessary cleanup
    account := pm.CleanupCacheAndFetchAccount(username)

    // If a previous payment existed, remove it
    if account.Payment != nil {
        account.Remove(account.Payment.Identifier)
    }

    // Set or update the payment details
    account.Payment = payment

    // Add or update the related Path entry with a new timestamp
    account.Add(payment.Identifier, PeerAccount{}, PeerAccount{})  // No PeerAccount details needed

    // Reinsert to manage any possible race condition, though very unlikely
    pm.reinsert(username, account)
}

// Wrapper for initiating an outgoing payment
func (pm *PathManager) InitiateOutgoingPayment(username, paymentID string, counterpart PeerAccount) {
    pm.initiatePayment(username, paymentID, 0, counterpart)  // 0 for outgoing
}

// Wrapper for initiating an incoming payment
func (pm: *PathManager) InitiateIncomingPayment(username, paymentID string, counterpart PeerAccount) {
    pm.initiatePayment(username, paymentID, 1, counterpart)  // 1 for incoming
}
