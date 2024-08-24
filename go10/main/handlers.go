package main

import (
    "ripple/types"
    "ripple/handlers/trustlines/client_trustlines"
    "ripple/handlers/trustlines/server_trustlines"
    "ripple/handlers/payments/client_payments"
    "ripple/handlers/payments/server_payments"
)

// CommandHandler defines the type for command handling functions
type CommandHandler func(session types.Session)

// CommandHandlers maps command bytes to handler functions
var commandHandlers = [256]CommandHandler{
    0:   client_trustlines.SetTrustline,     // Client Command
    1:   client_trustlines.SyncTrustlineIn,  // Client Command
    2:   client_trustlines.SyncTrustlineOut, // Client Command
    3:   client_trustlines.GetTrustlineIn,   // Client Command
    4:   client_trustlines.GetTrustlineOut,  // Client Command
    5:   client_payments.NewPaymentOut,      // Client Command
    6:   client_payments.NewPaymentIn,       // Client Command
    7:   client_payments.GetPayment,         // Client Command

    127: server_trustlines.SetTrustline,     // Server Command
    128: server_trustlines.GetTrustline,     // Server Command
    129: server_trustlines.SetSyncOut,       // Server Command
    130: server_trustlines.SetTimestamp,     // Server Command
    131: server_payments.FindPathOut,        // Server Command
    132: server_payments.FindPathIn    ,     // Server Command
    133: server_payments.PathRecurse  ,      // Server Command
    // Other indices are nil by default
}
