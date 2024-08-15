package pathfinding

import (
    "sync"
    "time"
)

// PathManager manages all Account entries in a system.
type PathManager struct {
    Accounts map[string]*Account // Map usernames to their respective Accounts.
    mu       sync.Mutex          // Protects the Accounts map.
}

// NewPathManager initializes and returns a new PathManager instance.
func NewPathManager() *PathManager {
    return &PathManager{
        Accounts: make(map[string]*Account), // Properly initialize the map.
    }
}

// Add creates a new account every time, overwriting any existing one.
func (pm *PathManager) Add(username string) *Account {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    // Create a new Account with the current time as the LastModified timestamp
    account := &Account{
        Username:     username,
        LastModified: time.Now(),
        Paths:        make(map[string]*Path),
    }

    // Always add the new account to the Accounts map, overwriting any existing one
    pm.Accounts[username] = account
    return account
}

// Find retrieves an account from the manager.
func (pm *PathManager) Find(username string) *Account {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    if account, exists := pm.Accounts[username]; exists {
        return account
    }
    return nil
}

// Remove deletes an account from the manager.
func (pm *PathManager) Remove(username string) {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    delete(pm.Accounts, username)
}

// Add creates and adds a new Path to an Account.
func (account *Account) Add(identifier string, incoming, outgoing PeerAccount) {
    newPath := NewPath(identifier, incoming, outgoing)
    account.Paths[identifier] = newPath
}

// Find retrieves a Path from an Account using the identifier.
func (account *Account) Find(identifier string) *Path {
    if path, exists := account.Paths[identifier]; exists {
        return path
    }
    return nil
}

// Remove deletes a Path from an Account using the identifier.
func (account *Account) Remove(identifier string) {
    delete(account.Paths, identifier)
}
