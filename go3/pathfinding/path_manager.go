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

// AddAccount adds a new account to the PathManager's linked list and returns the new AccountNode.
func (pm *PathManager) AddAccount(username string) *AccountNode {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    newNode := &AccountNode{
        Username:     username,
        LastModified: time.Now(),
        PathFinding:  nil, // Initialize with no pathfinding entries
        Next:         pm.head,
    }
    pm.head = newNode
    return newNode // Return the newly created AccountNode
}

// AddPathEntry adds a new PathEntry to the AccountNode's PathFinding linked list.
// It takes the incoming and outgoing PeerAccount, as well as a unique identifier.
func (node *AccountNode) AddPathEntry(identifier string, incoming, outgoing PeerAccount) {
    newEntry := &PathEntry{
        Identifier: identifier,
        Timestamp:  time.Now(),
        Incoming:   incoming,
        Outgoing:   outgoing,
        Next:      node.PathFinding, // Insert at the beginning
    }

    // Update the PathFinding list
    node.PathFinding = newEntry
}

// FindAccount searches for a specific account in the PathManager's linked list
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

// FindPathEntry checks if the given identifier exists in the PathFinding linked list,
// removes any expired entries based on the configured timeout duration,
// and returns the PathEntry for the identifier if it is found.
func (node *AccountNode) Find(identifier string) *PathEntry {
    // Use the BaseList's Find method to search for the PathEntry
    baseNode := node.BaseList.Find(identifier)

    // Simplified type assertion to PathEntry
    if baseNode != nil {
        return baseNode.(*PathEntry) // Directly return the asserted PathEntry
    }
    return nil // Not found or expired
}
