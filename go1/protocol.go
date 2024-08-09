// Datagram represents the updated structure
type Datagram struct {
    ClientOrServer   byte
    Username         [32]byte
    PeerUsername     [32]byte
    PeerServerAddress [32]byte
    Command          byte
    Arguments        [256]byte
    Counter          [4]byte
}

// CommandHandler defines the type for command handling functions
type CommandHandler func(session Session)

// CommandHandlers maps command bytes to handler functions
var commandHandlers = [256]CommandHandler{
    0x01: handleClientCommand1,
    0x02: handleClientCommand2,
    // Add more command handlers as needed
}
