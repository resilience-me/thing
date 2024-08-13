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

// Payment struct to hold details about an ongoing payment.
type Payment struct {
    Identifier string // Unique identifier for the payment
    InOrOut    bool   // True if outgoing (sender), false if incoming (receiver)
}

// AccountNode struct to represent each node's information in the account linked list.
type AccountNode struct {
    linkedlist.BaseNode           // Embedding the base struct for shared fields
    linkedlist.BaseList           // Embed BaseList to manage the linked list of BaseNodes
    Payment          *Payment     // Pointer to a Payment struct if a payment is active
}
