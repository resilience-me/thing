package pathfinding

import (
    "time"
    "ripple/config"
)

func (pm *PathManager) FetchAndRefresh(username string) *Account {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    if account, exists := pm.Accounts[username]; exists {
        newTimeout := time.Now().Add(config.PathFindingTimeout)
        // Ensure reinsert does not lower Timeout timer for an account currently committed to a payment
        if newTimeout.After(account.Timeout) {
            account.Timeout = newTimeout
        }
        return account
    }
    return nil
}
