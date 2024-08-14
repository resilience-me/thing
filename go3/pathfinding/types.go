package pathfinding

import (
    "sync"
    "time"
)

// PeerAccount holds the details about a peer account.
type PeerAccount struct {
    Username      string
    ServerAddress string
}

// PathNode is the replacement for PathEntry, adapted for use with a map.
type PathNode struct {
    Identifier   [32]byte // Using array for fixed-size identifiers.
    Timestamp    time.Time
    Incoming     PeerAccount
    Outgoing     PeerAccount
    CounterIn    int
    CounterOut   map[string]int // Map for outgoing counters by username.
}

// AccountNode holds all pathfinding related nodes and payment information.
type AccountNode struct {
    Username      string
    LastModified  time.Time
    Paths         map[string]*PathNode // Maps identifiers to PathNode.
    Payment       *Payment
}

// Payment structure adapted for use with AccountNode.
type Payment struct {
    Identifier string
    InOrOut    bool // True for outgoing, false for incoming.
}
