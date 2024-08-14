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

// AddAccount adds a new account or returns an existing one.
func (pm *PathManager) AddAccount(username string) *Account {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    // Check if the account already exists and return it if so
    if account, exists := pm.Accounts[username]; exists {
        return account
    }

    // If not exists, create a new Account with the current time as the LastModified timestamp
    account := &Account{
        Username:     username,
        LastModified: time.Now(), // Set the LastModified to the current time
        Paths:        make(map[string]*Path),
    }

    // Add the new account to the Accounts map
    pm.Accounts[username] = account
    return account
}

// FindAccount retrieves an account from the manager.
func (pm *PathManager) FindAccount(username string) *Account {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    if account, exists := pm.Accounts[username]; exists {
        return account
    }
    return nil
}

// RemoveAccount removes an account from the manager.
func (pm *PathManager) RemoveAccount(username string) {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    delete(pm.Accounts, username)
}

// AddPath adds a new Path to an Account.
func (account *Account) AddPath(identifier string, incoming, outgoing PeerAccount) {
    // Create a new Path entry
    newPath := &Path{
        Identifier:   identifier,
        Timestamp:    time.Now(),
        Incoming:     incoming,
        Outgoing:     outgoing,
        CounterIn:    0,
        CounterOut:   make(map[string]int),
    }
    // Add the new path to the Account's Paths map
    account.Paths[identifier] = newPath
}

// FindPath retrieves a Path from an Account.
func (account *Account) FindPath(identifier string) *Path {
    // Direct access to the path using the map
    if path, exists := account.Paths[identifier]; exists {
        return path
    }
    return nil
}

// RemovePath removes a Path from an Account.
func (account *Account) RemovePath(identifier string) {
    // Remove the path from the Account's Paths map
    delete(account.Paths, identifier)
}
