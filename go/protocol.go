package main

import (
    "net"
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

// Direct initialization of commandHandlers
var commandHandlers = [256]CommandHandler{
    handleSetTrustline, // Command 0
    handleGetTrustline, // Command 1
    // All other handlers are implicitly set to nil
}
