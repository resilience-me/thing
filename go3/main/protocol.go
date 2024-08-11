package main

package main

import (
    "ripple/handlers/client_trustlines"
    "ripple/handlers/server_trustlines"
)

// Datagram holds the structure of the incoming data
type Datagram struct {
    Command           byte
    Username          string
    PeerUsername      string
    PeerServerAddress string
    Arguments         [256]byte
    Counter           uint32       // Storing the counter as uint32
    Signature         [32]byte
}

// CommandHandler defines the type for command handling functions
type CommandHandler func(session Session)

// CommandHandlers maps command bytes to handler functions
var commandHandlers = [256]CommandHandler{
    0:   client_trustlines.SetTrustline,    // Client Command
    1:   client_trustlines.GetTrustlineOut, // Client Command
    2:   client_trustlines.GetTrustlineIn,  // Client Command
    128: client_trustlines.SetTrustline,    // Server Command
    129: server_trustlines.GetTrustline,    // Server Command
    130: server_trustlines.SetSyncOut,      // Server Command
    131: server_trustlines.SetTimestamp,    // Server Command
    // Other indices are nil by default
}

const (
    ClientTrustlines_SetTrustline      = 0
    ClientTrustlines_GetTrustlineOut   = 1
    ClientTrustlines_GetTrustlineIn    = 2
    ServerTrustlines_SetTrustline      = 128
    ServerTrustlines_GetTrustline      = 129
    ServerTrustlines_SetSyncOut        = 130
    ServerTrustlines_SetTimestamp      = 131
)
