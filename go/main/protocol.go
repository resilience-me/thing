package main

import (
    "net"
    "resilience/handlers/client"
    "resilience/handlers/server"
)

// Datagram holds the structure of the incoming data
type Datagram struct {
    Command        byte
    XUsername      [32]byte
    YUsername      [32]byte
    YServerAddress [32]byte
    Arguments      [256]byte
    Counter        [4]byte
    Signature      [32]byte
}

// CommandHandler defines the type for command handling functions
type CommandHandler func(Datagram, *net.UDPAddr)

// CommandHandlers holds the command handlers
var commandHandlers = [256]CommandHandler{
    0:   client.SetTrustline,    // Client Command
    128: server.SetTrustline,    // Server Command
    // Other indices are nil by default
}
