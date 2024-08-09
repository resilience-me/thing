// Will rewrite things to use tcp instead of udp, just to have the retransmission feature built-in. This is a quick first outline/sketch

package main

import (
    "fmt"
    "io"
    "net"
    "os"
)

// Datagram represents the structure you provided earlier
type Datagram struct {
    Command        byte
    XUsername      [32]byte
    YUsername      [32]byte
    YServerAddress [32]byte
    Arguments      [256]byte
    Counter        [4]byte
    Signature      [32]byte
}

// Define the BaseSession struct, embedding the Datagram
type BaseSession struct {
    Datagram
}

// Define the Session interface with a GetDatagram method
type Session interface {
    GetDatagram() *Datagram
}

// Implement the GetDatagram method for BaseSession
func (bs *BaseSession) GetDatagram() *Datagram {
    return &bs.Datagram
}

// Define the ClientSession struct
type ClientSession struct {
    BaseSession
    Conn net.Conn
}

// Define the ServerSession struct
type ServerSession struct {
    BaseSession
}

// SessionManager manages the processing of sessions, including client and server datagrams
type SessionManager struct {
    sessionCh      chan Session                // Create a channel for Session interfaces
    closedCh       chan [32]byte               // Channel for closed sessions
    activeHandlers map[[32]byte]bool           // Tracks active handlers by username
    queues         map[[32]byte][]Session      // Queues for sessions waiting to be processed
}

// HandlerContext holds the common parameters for handler functions
type HandlerContext struct {
    Session Session          // The session, which can be ClientSession or ServerSession
    CloseCh chan [32]byte    // Channel to signal completion
}

// CommandHandler defines the type for command handling functions
type CommandHandler func(ctx HandlerContext)

// CommandHandlers maps command bytes to handler functions
var commandHandlers = [256]CommandHandler{
    0x01: handleClientCommand1,
    0x02: handleClientCommand2,
    // Add more command handlers as needed
}

func (m *SessionManager) run() {
    for {
        select {
        case session := <-m.sessionCh:
            username := session.GetDatagram().XUsername

            if !m.processors[username] {
                // No active handler, create HandlerContext and start processing
                m.processors[username] = true

                go m.handleSession(session)
            } else {
                // Processor is active, enqueue the session
                m.queues[username] = append(m.queues[username], session)
            }

        case username := <-m.closedCh:
            // Processor finished, check if there are queued sessions
            if queue, exists := m.queues[username]; exists && len(queue) > 0 {
                // Start a new processor with the next session
                nextSession := queue[0]
                m.queues[username] = queue[1:]

                go m.handleSession(nextSession)
            } else {
                // No sessions left, mark processor as not running
                delete(m.activeHandlers, username)
            }
        }
    }
}

// handleSession creates a new context and processes the datagram
func (m *SessionManager) handleSession(session Session) {
    defer func() {
        m.closedCh <- session.GetDatagram().XUsername
    }()
    
    command := session.GetDatagram().Command

    // Look up the command handler
    handler := commandHandlers[command]
    if handler == nil {
        fmt.Printf("Unknown command: %d\n", command)
        return
    }
    
    // Execute the handler
    handler(session)
}

// bytesToDatagram populates a Datagram struct from a byte slice
func bytesToDatagram(dg *Datagram, buf []byte) {
    dg.Command = buf[0]
    copy(dg.XUsername[:], buf[1:33])
    copy(dg.YUsername[:], buf[33:65])
    copy(dg.YServerAddress[:], buf[65:97])
    copy(dg.Arguments[:], buf[97:353])
    copy(dg.Counter[:], buf[353:357])
    copy(dg.Signature[:], buf[357:389])
}

// handleConnection reads datagrams from the connection and sends them to the SessionManager
func (m *SessionManager) handleConnection(conn net.Conn) {
    buf := make([]byte, 389) // Create a buffer with the size of the Datagram

    // Read the full datagram into the buffer
    _, err := io.ReadFull(conn, buf)
    if err != nil {
        if err == io.EOF {
            fmt.Println("Connection closed by client")
        } else {
            fmt.Printf("Error reading datagram: %v\n", err)
        }
        return
    }

    // Ensure the buffer is the correct size
    if len(buf) < 389 {
        fmt.Println("Buffer is too small")
        return
    }

    // Determine if this is a server or client command based on the command byte
    isServerCommand := (buf[0] & 0x80) != 0 // Check the MSB of the Command byte

    if isServerCommand {
        // Create and populate a ServerSession
        serverSession := &ServerSession{}
        bytesToDatagram(&serverSession.Datagram, buf)
        m.sessionCh <- serverSession
        conn.Close() // Close the connection for server sessions
    } else {
        // Create and populate a ClientSession
        clientSession := &ClientSession{
            Conn: conn,
        }
        bytesToDatagram(&clientSession.Datagram, buf)
        m.sessionCh <- clientSession
        // Connection remains open for client sessions
    }
}

// Main function with inlined server logic
func main() {
    manager := &SessionManager{
        sessionCh:      make(chan Session),
        closedCh:       make(chan [32]byte),
        activeHandlers: make(map[[32]byte]bool),
        queues:         make(map[[32]byte][]Session),
    }
    go manager.run()

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
        go manager.handleConnection(conn)
    }
}
