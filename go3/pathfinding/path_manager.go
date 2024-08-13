package main

import (
    "crypto/sha256"
    "fmt"
    "time"
)

type AccountNode struct {
    Username     string
    LastModified time.Time
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

type PathManager struct {
    head *AccountNode
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

// FindAccount searches for a node by username
func (pm *PathManager) FindAccount(username string) *AccountNode {
    current := pm.head
    for current != nil {
        if current.Username == username {
            return current
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
