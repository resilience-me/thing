package pathfinding

import (
  "ripple/linkedlist"
)

// PeerAccount represents a peer in the network.
type PeerAccount struct {
    Username      string
    ServerAddress string
}

// PathEntry represents an entry in the pathfinding linked list.
type PathEntry struct {
    linkedlist.BaseNode // Embedding the base struct for shared fields
    Incoming           PeerAccount
    Outgoing           PeerAccount
    CounterIn          int
    CounterOut         map[string]int
}

// AccountNode represents a node in the account linked list.
type AccountNode struct {
    linkedlist.BaseNode           // Embedding the base struct for shared fields
    PathFinding        *PathEntry // Linked list of PathEntry nodes
}
