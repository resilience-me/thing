package main

import (
    "sync"
    "time"
)

// PeerAccount holds details about a peer account.
type PeerAccount struct {
    Username      string
    ServerAddress string
}

// Path replaces PathNode, tailored for use with a map and string identifiers.
type Path struct {
    Identifier   string          // Identifier for the path.
    Timestamp    time.Time       // Timestamp of the last update.
    Incoming     PeerAccount     // Details of the incoming peer.
    Outgoing     PeerAccount     // Details of the outgoing peer.
    CounterIn    int             // Counter for incoming paths.
    CounterOut   map[string]int  // Map for outgoing counters by username.
}

// Account holds all path-related information and payment details.
type Account struct {
    Username      string
    LastModified  time.Time
    Paths         map[string]*Path // Maps string identifiers to Path.
    Payment       *Payment
}

// Payment structure adapted for use with Account.
type Payment struct {
    Identifier string
    InOrOut    bool // True for outgoing, false for incoming.
}
