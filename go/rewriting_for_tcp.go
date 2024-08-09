// Will rewrite things to use tcp instead of udp, just to have the retransmission feature built-in. This is a quick first outline/sketch

package main

import (
    "fmt"
    "io"
    "net"
    "os"
)

// Datagram represents the structure you provided earlier.
type Datagram struct {
    Command        byte
    XUsername      [32]byte
    YUsername      [32]byte
    YServerAddress [32]byte
    Arguments      [256]byte
    Counter        [4]byte
    Signature      [32]byte
}

// HandlerContext holds the common parameters for handler functions
type HandlerContext struct {
    Datagram *Datagram    // Pointer to Datagram
    Conn     net.Conn     // TCP connection, for client-server responses
    CloseCh  chan [32]byte  // Channel to signal completion
}

// CommandHandler defines the type for command handling functions
type CommandHandler func(ctx HandlerContext)

// CommandHandlers maps command bytes to handler functions
var commandHandlers = [256]CommandHandler{
    0x01: handleClientCommand1,
    0x02: handleClientCommand2,
    // Add more command handlers as needed
}

// AccountManager manages the processing of datagrams per account
type AccountManager struct {
    datagramCh chan Datagram              // Channel for incoming datagrams
    closedCh   chan [32]byte              // Channel for signals from processors
    processors map[[32]byte]bool          // Active processors
    queues     map[[32]byte][]Datagram    // Queues for pending datagrams per account
}

// NewAccountManager creates a new AccountManager
func NewAccountManager() *AccountManager {
    return &AccountManager{
        datagramCh: make(chan Datagram),
        closedCh:   make(chan [32]byte),
        processors: make(map[[32]byte]bool),
        queues:     make(map[[32]byte][]Datagram),
    }
}

// Run listens for datagrams and manages their processing
func (m *AccountManager) Run() {
    for {
        select {
        case data := <-m.datagramCh:
            username := data.XUsername
            if !m.processors[username] {
                // No active processor, start one
                m.processors[username] = true
                go m.ProcessDatagram(data, nil, m.closedCh)
            } else {
                // Processor is active, enqueue the datagram
                m.queues[username] = append(m.queues[username], data)
            }

        case username := <-m.closedCh:
            // Processor finished, check if there are queued datagrams
            if queue, exists := m.queues[username]; exists && len(queue) > 0 {
                // Start a new processor with the next datagram
                nextDg := queue[0]
                m.queues[username] = queue[1:]
                go m.ProcessDatagram(nextDg, nil, m.closedCh)
            } else {
                // No datagrams left, mark processor as not running
                m.processors[username] = false
            }
        }
    }
}

// ProcessDatagram creates a new context and processes the datagram
func (m *AccountManager) ProcessDatagram(datagram Datagram, conn net.Conn, closeCh chan [32]byte) {
    ctx := HandlerContext{
        Datagram: &datagram,
        Conn:     conn,
        CloseCh:  closeCh,
    }

    // Look up the command handler
    handler := commandHandlers[ctx.Datagram.Command]
    if handler == nil {
        fmt.Printf("Unknown command: %d\n", ctx.Datagram.Command)
        ctx.CloseCh <- ctx.Datagram.XUsername
        return
    }

    // Execute the handler
    handler(ctx)
}

// datagramBytes provides a slice that covers the entire datagram for reading
func datagramBytes(d *Datagram) []byte {
    size := 1 + 32 + 32 + 32 + 256 + 4 + 32 // Total size of the Datagram struct
    return (*[389]byte)(d)[:size]
}

// Example client-server command handler 1
func handleClientCommand1(ctx HandlerContext) {
    defer func() {
        ctx.CloseCh <- ctx.Datagram.XUsername
    }()

    fmt.Println("Handling Client Command 1")
    // The handler may or may not send a response
}

// Example client-server command handler 2
func handleClientCommand2(ctx HandlerContext) {
    defer func() {
        ctx.CloseCh <- ctx.Datagram.XUsername
    }()

    fmt.Println("Handling Client Command 2")
    // The handler may or may not send a response
}

// isServerCommand checks if the command indicates a server command
func isServerCommand(command byte) bool {
    return (command & 0x80) != 0 // 0x80 is 10000000 in binary
}

// handleConnection reads datagrams from the connection and sends them to the AccountManager
func handleConnection(conn net.Conn, manager *AccountManager) {
    var datagram Datagram
    err := io.ReadFull(conn, datagramBytes(&datagram))
    if err != nil {
        if err == io.EOF {
            fmt.Println("Connection closed by client")
        } else {
            fmt.Printf("Error reading datagram: %v\n", err)
        }
        return
    }

    // Use the isServerCommand function to determine session type
    if isServerCommand(datagram.Command) {
        // Handle server command: Send to server channel
        manager.serverCh <- ServerSession{Datagram: datagram}
        // Close the connection after ensuring the ServerSession has been sent
        conn.Close()
    } else {
        // Handle client command: Keep the connection open for potential responses
        manager.clientCh <- ClientSession{Datagram: datagram, Conn: conn}
        // Connection will remain open for further processing
    }
}

// Main function with inlined server logic
func main() {
    manager := NewAccountManager()
    go manager.Run()

    // Start the TCP server on port 2012
    listener, err := net.Listen("tcp", ":2012")
    if err != nil {
        fmt.Printf("Error starting TCP server: %v\n", err)
        os.Exit(1)
    }
    defer listener.Close()

    fmt.Println("Listening on port 2012...")

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Printf("Error accepting connection: %v\n", err)
            continue
        }

        // Handle each connection in a separate goroutine
        go handleConnection(conn, manager)
    }
}
