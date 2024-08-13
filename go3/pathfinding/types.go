package pathfinding

import (
  "ripple/pathfinding/linkedlist"
)

// PeerAccount represents a peer in the network.
type PeerAccount struct {
    Username      string
    ServerAddress string
}

// PathNode represents an entry in the pathfinding linked list.
type PathNode struct {
    linkedlist.BaseNode // Embedding the base struct for shared fields
    Incoming           PeerAccount
    Outgoing           PeerAccount
    CounterIn          int
    CounterOut         map[string]int
}

// AccountNode represents a node in the account linked list.
type AccountNode struct {
    linkedlist.BaseNode           // Embedding the base struct for shared fields
    linkedlist.BaseList // Embed BaseList to manage the linked list of BaseNodes
}
