package pathfinding

import (
    "sync"
)

// PathManager manages all AccountNodes in a system.
type PathManager struct {
    Accounts map[string]*AccountNode // Map usernames to their respective AccountNodes.
    mu       sync.Mutex // Protects the Accounts map.
}

// NewPathManager initializes and returns a new PathManager instance
func NewPathManager() *PathManager {
    return &PathManager{}
}

// AddAccount adds a new account or returns an existing one.
func (pm *PathManager) AddAccount(username string) *AccountNode {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    if node, exists := pm.Accounts[username]; exists {
        return node
    }

    node := &AccountNode{
        Username: username,
        Paths:    make(map[string]*PathNode),
    }
    pm.Accounts[username] = node
    return node
}

// FindAccount retrieves an account from the manager.
func (pm *PathManager) FindAccount(username string) *AccountNode {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    if node, exists := pm.Accounts[username]; exists {
        return node
    }
    return nil
}

// RemoveAccount removes an account from the manager.
func (pm *PathManager) RemoveAccount(username string) {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    delete(pm.Accounts, username)
}

// Add adds a new PathNode to the AccountNode's PathFinding linked list.
func (node *AccountNode) Add(identifier string, incoming, outgoing PeerAccount) {
    newEntry := &PathNode{
        BaseNode: linkedlist.BaseNode{Identifier: identifier},
        Incoming: incoming,
        Outgoing: outgoing,
    }
    node.BaseList.Add(&newEntry.BaseNode)
}

// Find checks if the given identifier exists in the PathFinding linked list,
// removes any expired entries based on the configured timeout duration,
// and returns the PathNode for the identifier if it is found.
func (node *AccountNode) Find(identifier string) *PathNode {
    baseNode := node.BaseList.Find(identifier)

    if baseNode != nil {
        return baseNode.(*PathNode)
    }
    return nil
}

func (node *AccountNode) Remove(identifier string) {
    node.PathList.Remove(identifier)
}
