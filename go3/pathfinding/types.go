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

// Account holds all path-related information and payment details
type Account struct {
    Username      string
    Cleanup       time.Time
    Paths         map[string]*Path // Maps string identifiers to Path.
    Payment       *Payment
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
