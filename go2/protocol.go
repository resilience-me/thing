package main

// Datagram represents the updated structure
type Transaction struct {
    Command            byte
    Username           string
    PeerUsername       string
    PeerServerAddress  string
    Arguments          [256]byte
    Counter            [4]byte
}

// CommandHandler defines the type for command handling functions
type CommandHandler func(session Session)

// CommandHandlers maps command bytes to handler functions
var commandHandlers = [256]CommandHandler{
    0x01: handleClientCommand1,
    0x02: handleClientCommand2,
    // Add more command handlers as needed
}
