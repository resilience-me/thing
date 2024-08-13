package pathfinding

import (
    "crypto/sha256"
    "fmt"
    "time"
    "ripple/config"
)

// AccountNode represents a node in the linked list
type AccountNode struct {
    Username     string
    LastModified time.Time
    PathFinding  *PathEntry // Linked list of PathEntry nodes
    Next         *AccountNode
}

type PeerAccount struct {
    Username      string
    ServerAddress string
}

type PathEntry struct {
    Identifier [32]byte
    Timestamp  time.Time
    Incoming   PeerAccount
    Outgoing   PeerAccount
    Next       *PathEntry
}

// PathManager manages the linked list of accounts
type PathManager struct {
    head *AccountNode
    mu   sync.Mutex // Mutex to protect access to the linked list
}

// NewPathManager creates a new PathManager with an initial list of accounts
func NewPathManager() *PathManager {
    var head *AccountNode

    for _, username := range initialAccounts {
        head = &AccountNode{
            Username:     username,
            LastModified: time.Now(),
            Next:         head,
        }
    }

    return &PathManager{
        head: head,
    }
}

// AddAccount adds a new account to the beginning of the list
func (pm *PathManager) AddAccount(username string) {
    newNode := &AccountNode{
        Username:     username,
        LastModified: time.Now(),
        Next:         pm.head,
    }
    pm.head = newNode
}

// FindAccount searches for a specific account in the PathManager's linked list.
// If the account is found but has timed out, it is removed from the list, and nil is returned.
// If the account is found and has not timed out, it is returned.
// Along the way, any other accounts that have timed out are also removed from the list.
// The function returns nil if the account is not found or has been removed due to timeout.
func (pm *PathManager) FindAccount(username string) *AccountNode {
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


// DisplayAccounts prints out the entire linked list
func (pm *PathManager) DisplayAccounts() {
    current := pm.head
    for current != nil {
        fmt.Printf("Username: %s, Last Modified: %s\n", current.Username, current.LastModified)
        current = current.Next
    }
}

// HandlePathRequest processes a single hop for either incoming or outgoing path
func (pm *PathManager) HandlePathRequest(identifier [32]byte, isOutgoing bool, requestOrigin PeerAccount) *PathEntry {
    return handlePathRequest(pm.head, identifier, isOutgoing, requestOrigin)
}
