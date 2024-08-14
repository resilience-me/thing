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

// Add adds a new account or returns an existing one.
func (pm *PathManager) Add(username string) *Account {
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
    newPath := &Path{
        Identifier:   identifier,
        Timestamp:    time.Now(),
        Incoming:     incoming,
        Outgoing:     outgoing,
        CounterIn:    0,
        CounterOut:   make(map[string]int),
    }
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
