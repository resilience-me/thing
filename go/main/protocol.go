package main

import (
    "net"
    "resilience/handlers/client_trustlines"
    "resilience/handlers/server_trustlines"
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
    Datagram *main.Datagram // Pointer to Datagram
    Addr     *net.UDPAddr
    Conn     *net.UDPConn
}

// CommandHandler defines the type for command handling functions
type CommandHandler func(HandlerContext)

// CommandHandlers holds the command handlers
var commandHandlers = [256]CommandHandler{
    0:   client_trustlines.SetTrustline,    // Client Command
    1:   client_trustlines.GetTrustline,    // Client Command
    128: client_trustlines.SetTrustline,    // Server Command
    129: server_trustlines.GetTrustline,    // Server Command
    130: server_trustlines.SetSyncOut,      // Server Command
    131: server_trustlines.SetTimestamp,    // Server Command
    // Other indices are nil by default
}

const (
    ClientTrustlines_SetTrustline   = 0
    ClientTrustlines_GetTrustline   = 1
    ServerTrustlines_SetTrustline   = 128
    ServerTrustlines_GetTrustline   = 129
    ServerTrustlines_SetSyncOut     = 130
    ServerTrustlines_SetTimestamp   = 131
)
