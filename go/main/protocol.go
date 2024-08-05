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

type ResponseDatagram {
    Nonce        [32]byte
    Result       [256]byte
    Signature    [32]byte
}

// HandlerContext holds the common parameters for handler functions
type HandlerContext struct {
    Datagram main.Datagram
    Addr     *net.UDPAddr
    Conn     *net.UDPConn
}

// CommandHandler defines the type for command handling functions
type CommandHandler func(HandlerContext)

// CommandHandlers holds the command handlers
var commandHandlers = [256]CommandHandler{
    0:   client.SetTrustline,    // Client Command
    128: server.SetTrustline,    // Server Command
    129: server.SetSyncCounter,  // Server Command
    130: server.GetTrustline,    // Server Command
    // Other indices are nil by default
}

const (
    Client_SetTrustline   = 0
    Server_SetTrustline   = 128
    Server_SetSyncCounter = 129
    Server_GetTrustline   = 130
)
