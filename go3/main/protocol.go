package main

package main

import (
    "ripple/handlers/trustlines/client_trustlines"
    "ripple/handlers/trustlines/server_trustlines"
    "ripple/handlers/payments"
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
    0:   client_trustlines.SetTrustline,     // Client Command
    1:   client_trustlines.SyncTrustlineIn,  // Client Command
    2:   client_trustlines.SyncTrustlineOut, // Client Command
    3:   client_trustlines.GetTrustlineIn,   // Client Command
    4:   client_trustlines.GetTrustlineOut,  // Client Command
    127: server_trustlines.SetTrustline,     // Server Command
    128: server_trustlines.GetTrustline,     // Server Command
    129: server_trustlines.SetSyncOut,       // Server Command
    130: server_trustlines.SetTimestamp,     // Server Command
    131: payments.NewPaymentOut,             // Server Command
    132: payments.NewPaymentIn,              // Server Command
    133: payments.NewPathFindingOut,         // Server Command
    134: payments.PathFindingOut,            // Server Command
    135: payments.PathFindingIn,             // Server Command
    136: payments.PathFindingRecurse,        // Server Command
    // Other indices are nil by default
}

const (
    ClientTrustlines_SetTrustline      = 0
    ClientTrustlines_SyncTrustlineIn   = 1
    ClientTrustlines_SyncTrustlineOut  = 2
    ClientTrustlines_GetTrustlineIn    = 3
    ClientTrustlines_GetTrustlineOut   = 4
    ServerTrustlines_SetTrustline      = 127
    ServerTrustlines_GetTrustline      = 128
    ServerTrustlines_SetSyncOut        = 129
    ServerTrustlines_SetTimestamp      = 130
    Payments_NewPaymentOut             = 131
    Payments_NewPaymentIn              = 132
    Payments_NewPathFindingOut         = 133
    Payments_PathFindingOut            = 134
    Payments_PathFindingIn             = 135
    Payments_PathFindingRecurse        = 136
)
