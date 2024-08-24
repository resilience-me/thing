package pathfinding

import "time"

func FetchAndRefresh() *Account {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    if account, exists := pm.Accounts[username]; exists {
        newCleanup := time.Now().Add(config.PathFindingTimeout)
        // Ensure reinsert does not lower Cleanup timer for an account currently committed to a payment
        if newCleanup.After(account.Cleanup) {
            account.Cleanup = newCleanup
        }
        return account
    }
    return nil
}
