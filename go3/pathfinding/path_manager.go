package pathfinding

import (
    "crypto/sha256"
    "fmt"
    "time"
    "ripple/config"
)

type PeerAccount struct {
    Username      string
    ServerAddress string
}

// PathEntry represents an entry in the pathfinding linked list
type PathEntry struct {
    Identifier [32]byte
    Timestamp  time.Time
    Incoming   PeerAccount
    Outgoing   PeerAccount
    Counter    uint32
    Next       *PathEntry
}

// AccountNode represents a node in the linked list
type AccountNode struct {
    Username     string
    LastModified time.Time
    PathFinding  *PathEntry // Linked list of PathEntry nodes
    Next         *AccountNode
}

// PathManager manages the linked list of accounts
type PathManager struct {
    head *AccountNode
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

// FindAccount searches for a specific account in the PathManager's linked list
// and returns it if found. This method also removes any accounts that have
// timed out (based on the LastModified timestamp) as it traverses the list.
// Thread safety is ensured using a mutex.
func (pm *PathManager) FindAccount(username string) *AccountNode {
    pm.mu.Lock()         // Lock the mutex before accessing shared data
    defer pm.mu.Unlock() // Ensure the mutex is unlocked when the function returns

    var prev *AccountNode
    current := pm.head
    now := time.Now()

    for current != nil {
        isTarget := current.Username == username

        if now.Sub(current.LastModified) > config.PathFindingTimeout {
            // Remove timed-out node, whether it's the target or not
            if prev == nil {
                pm.head = current.Next
            } else {
                prev.Next = current.Next
            }

            // If the timed-out node was the target, return nil immediately
            if isTarget {
                return nil
            }
        } else {
            // If it's not timed out, check if it's the target
            if isTarget {
                return current // Target found and not timed out
            }
            prev = current
        }

        current = current.Next
    }
    return nil
}

// AddPathEntry adds a new PathEntry to the AccountNode's PathFinding linked list.
// It takes the incoming and outgoing PeerAccount, as well as a unique identifier.
func (node *AccountNode) AddPathEntry(identifier [32]byte, incoming, outgoing PeerAccount) {
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

// FindPathEntry checks if the given identifier exists in the PathFinding linked list,
// removes any expired entries based on the configured timeout duration,
// and returns the PathEntry for the identifier if it is found.
func (node *AccountNode) FindPathEntry(identifier [32]byte) *PathEntry {
    now := time.Now()
    var prev *PathEntry
    current := node.PathFinding

    for current != nil {
        isTarget := current.Identifier == identifier

        if now.Sub(current.Timestamp) > config.PathFindingTimeout {
            // Remove expired entry
            if prev == nil {
                node.PathFinding = current.Next
            } else {
                prev.Next = current.Next
            }
            // If the timed-out node was the target, return nil immediately
            if isTarget {
                return nil
            }

        } else {
            // If it's not timed out, check if it's the target
            if isTarget {
                return current // Target found and not timed out
            }
            prev = current
        }

        current = current.Next
    }

    return nil
}

// HandlePathRequest processes a single hop for either incoming or outgoing path
func (pm *PathManager) HandlePathRequest(identifier [32]byte, isOutgoing bool, requestOrigin PeerAccount) *PathEntry {
    return handlePathRequest(pm.head, identifier, isOutgoing, requestOrigin)
}
