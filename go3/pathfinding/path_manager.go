package pathfinding

// Add creates a new account every time, overwriting any existing one.
func (pm *PathManager) Add(username string) *Account {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    account := NewAccount(username)
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
func (account *Account) Add(identifier string, amount uint32, incoming, outgoing PeerAccount) {
    newPath := NewPath(identifier, amount, incoming, outgoing)
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
