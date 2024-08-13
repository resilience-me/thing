package main

package main

import (
    "ripple/handlers/trustlines/client_trustlines"
    "ripple/handlers/trustlines/server_trustlines"
    "ripple/handlers/payments/client_payments"
    "ripple/handlers/payments/server_payments"
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
    5:   client_payments.NewPaymentOut,      // Client Command
    6:   client_payments.NewPaymentIn,       // Client Command
    127: server_trustlines.SetTrustline,     // Server Command
    128: server_trustlines.GetTrustline,     // Server Command
    129: server_trustlines.SetSyncOut,       // Server Command
    130: server_trustlines.SetTimestamp,     // Server Command
    131: server_payments.NewPathFindingOut,  // Server Command
    132: server_payments.PathFindingOut,     // Server Command
    133: server_payments.PathFindingIn,      // Server Command
    134: server_payments.PathFindingRecurse, // Server Command
    // Other indices are nil by default
}

const (
    ClientTrustlines_SetTrustline      = 0
    ClientTrustlines_SyncTrustlineIn   = 1
    ClientTrustlines_SyncTrustlineOut  = 2
    ClientTrustlines_GetTrustlineIn    = 3
    ClientTrustlines_GetTrustlineOut   = 4
    ClientPayments_NewPaymentOut       = 5
    ClientPayments_NewPaymentIn        = 6
    ServerTrustlines_SetTrustline      = 127
    ServerTrustlines_GetTrustline      = 128
    ServerTrustlines_SetSyncOut        = 129
    ServerTrustlines_SetTimestamp      = 130
    ServerPayments_NewPathFindingOut   = 131
    ServerPayments_PathFindingOut      = 132
    ServerPayments_PathFindingIn       = 133
    ServerPayments_PathFindingRecurse  = 134
)
