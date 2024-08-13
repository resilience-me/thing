package pathfinding

import (
    "crypto/sha256"
    "fmt"
    "time"
    "ripple/config"
    "ripple/linkedlist"
)

// PathManager manages the linked list of accounts
type PathManager struct {
    linkedlist.BaseList
    mu   sync.Mutex // Mutex to protect access to the linked list
}

// NewPathManager initializes and returns a new PathManager instance
func NewPathManager() *PathManager {
    return &PathManager{}
}

// Add adds a new account to the PathManager's linked list and returns the new AccountNode.
func (pm *PathManager) Add(username string) *AccountNode {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    newNode := &AccountNode{
        Identifier: username,
    }

    pm.BaseList.Add(&newNode.BaseNode)

    return newNode // Return the newly created AccountNode
}

// Add adds a new PathNode to the AccountNode's PathFinding linked list.
func (node *AccountNode) Add(identifier string, incoming, outgoing PeerAccount) {
    newEntry := &PathNode{
        Identifier: identifier,
        Incoming:   incoming,
        Outgoing:   outgoing,
    }

    // Use BaseList's Add method to insert the new entry
    node.PathFinding.Add(&newEntry.BaseNode)
}

// Find searches for a specific account in the PathManager's linked list
// and returns it if found. Thread safety is ensured using a mutex.
func (pm *PathManager) Find(username string) *AccountNode {
    pm.mu.Lock()         // Lock the mutex before accessing shared data
    defer pm.mu.Unlock() // Ensure the mutex is unlocked when the function returns

    // Use the BaseList's Find method to search for the node
    baseNode := pm.BaseList.Find(username)

    // Simplified type assertion to AccountNode
    if baseNode != nil {
        return baseNode.(*AccountNode) // Directly return the asserted AccountNode
    }
    return nil // Not found or expired
}

// Find checks if the given identifier exists in the PathFinding linked list,
// removes any expired entries based on the configured timeout duration,
// and returns the PathNode for the identifier if it is found.
func (node *AccountNode) Find(identifier string) *PathNode {
    // Use the BaseList's Find method to search for the PathNode
    baseNode := node.BaseList.Find(identifier)

    // Simplified type assertion to PathNode
    if baseNode != nil {
        return baseNode.(*PathNode) // Directly return the asserted PathNode
    }
    return nil // Not found or expired
}
