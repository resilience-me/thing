package pathfinding

import (
    "time"
)

// PeerAccount holds details about a peer account
type PeerAccount struct {
    Username      string
    ServerAddress string
}

// NewPeerAccount is a constructor for creating a PeerAccount struct based on the provided username and server address.
func NewPeerAccount(username, serverAddress string) PeerAccount {
    return PeerAccount{
        Username:      username,
        ServerAddress: serverAddress,
    }
}

// Path replaces PathNode, tailored for use with a map and string identifiers
type Path struct {
    Identifier   string          // Identifier for the path
    Timeout      time.Time       // Direct expiration time for the path
    Amount       uint32
    Incoming     PeerAccount     // Details of the incoming peer
    Outgoing     PeerAccount     // Details of the outgoing peer
    Commit       bool
    CounterIn    uint32          // Counter for incoming paths
    CounterOut   map[string]int  // Map for outgoing counters by username
}

// NewPath is a constructor for creating a Path struct based on an identifier, incoming and outgoing PeerAccount, and amount.
func NewPath(identifier string, amount uint32, incoming, outgoing PeerAccount) *Path {
    return &Path{
        Identifier:   identifier,
        Timeout:      time.Now().Add(config.PathFindingTimeout), // Set the Timeout using PathFindingTimeout
        Amount:       amount,                                   // Set the amount
        Incoming:     incoming,
        Outgoing:     outgoing,
        CounterIn:    0,
        CounterOut:   make(map[string]int),
    }
}

// Account holds all path-related information and payment details
type Account struct {
    Username      string
    Cleanup       time.Time
    Paths         map[string]*Path // Maps string identifiers to Path.
    Payment       *Payment
}

// NewAccount creates and returns a new Account with the provided username.
func NewAccount(username string) *Account {
    return &Account{
        Username: username,
        Cleanup:  time.Now().Add(config.PathFindingTimeout), // Set the initial Cleanup time
        Paths:    make(map[string]*Path),
    }
}

// Payment structure adapted for use with Account
type Payment struct {
    Identifier  string
    Counterpart PeerAccount
    InOrOut     byte  // 0 for incoming, 1 for outgoing, stored as a single byte
}

// NewPayment is a constructor for creating a Payment struct based on an identifier, datagram, and inOrOut value.
func NewPayment(datagram *Datagram, identifier string, inOrOut byte) *Payment {
    // Initialize and return the Payment struct, using NewPeerAccount for the Counterpart field
    return &Payment{
        Identifier: identifier,
        Counterpart: NewPeerAccount(
            datagram.PeerUsername,
            datagram.PeerServerAddress,
        ),
        InOrOut: inOrOut,
    }
}
