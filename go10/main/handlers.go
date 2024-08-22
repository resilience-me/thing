package main

// CommandHandler defines the type for command handling functions
type CommandHandler func(session Session)

// CommandHandlers maps command bytes to handler functions
var commandHandlers = [256]CommandHandler{
    0:   client_trustlines.SetTrustline,     // Client Command
    1:   client_trustlines.SyncTrustlineIn,  // Client Command
    2:   client_trustlines.SyncTrustlineOut, // Client Command
    3:   client_trustlines.GetTrustlineIn,   // Client Command
    4:   client_trustlines.GetTrustlineOut,  // Client Command
    // Other indices are nil by default
}
