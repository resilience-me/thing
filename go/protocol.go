package main

import (
    "net"
    "resilience/handlers"
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
    handlers.setTrustline, // Command 0
    // All other handlers are implicitly set to nil
}
