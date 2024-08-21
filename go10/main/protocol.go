package main

// CommandHandler defines the type for command handling functions
type CommandHandler func(session Session)

// CommandHandlers maps command bytes to handler functions
var commandHandlers = [256]CommandHandler{
    // Other indices are nil by default
}
