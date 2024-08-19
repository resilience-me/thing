package main

// Datagram holds the structure of the incoming data
type Datagram struct {
    Command           byte
    Username          string
    PeerUsername      string
    PeerServerAddress string
    Arguments         [256]byte
    Counter           uint32
    Signature         [32]byte
}

// CommandHandler defines the type for command handling functions
type CommandHandler func(session Session)

// CommandHandlers maps command bytes to handler functions
var commandHandlers = [256]CommandHandler{
    // Other indices are nil by default
}
