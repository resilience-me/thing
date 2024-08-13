package pathfinding

import (
    "sync"
    "ripple/pathfinding/linkedlist"
)

// PathManager manages the linked list of accounts
type PathManager struct {
    linkedlist.BaseList
    mu sync.Mutex // Mutex to protect access to the linked list
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
        BaseNode: linkedlist.BaseNode{Identifier: username},
    }

    pm.BaseList.Add(&newNode.BaseNode)

    return newNode
}

// Find searches for a specific account in the PathManager's linked list and returns it if found.
func (pm *PathManager) Find(username string) *AccountNode {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    baseNode := pm.BaseList.Find(username)

    if baseNode != nil {
        return baseNode.(*AccountNode)
    }
    return nil
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
