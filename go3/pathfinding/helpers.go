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

// Reinsert updates the Cleanup field and reinserts the account if it was removed.
func (pm *PathManager) reinsert(username string, account *Account) {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    // Calculate the new proposed Cleanup time
    newCleanup := time.Now().Add(config.PathFindingTimeout)

    // Only update the Cleanup field if the new timeout is later than the current one
    if newCleanup.After(account.Cleanup) {
        account.Cleanup = newCleanup
    }

    // Reinsert the account
    pm.Accounts[username] = account
}


// InitiatePayment sets up or updates payment details for an account, creating the account if necessary.
func (pm *PathManager) InitiatePayment(username string, payment *Payment) {
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

    // Reinsert to manage a very unlikely race condition
    pm.reinsert(username, account)
}
