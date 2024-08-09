package main

import (
    "fmt"
    "io"
    "net"
    "os"
)

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

func (m *SessionManager) run() {
    for {
        select {
        case session := <-m.sessionCh:
            username := session.GetDatagram().Username

            if !m.activeHandlers[username] {
                // No active handler, create HandlerContext and start processing
                m.activeHandlers[username] = true

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

// handleSession manages the processing of a session's datagram, 
// ensuring the appropriate command handler is called and
// sending a closure signal when processing is complete.
func (m *SessionManager) handleSession(session Session) {
    defer func() {
        m.closedCh <- session.GetDatagram().Username
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
    dg.ClientOrServer = buf[0]
    copy(dg.Username[:], buf[1:33])
    copy(dg.PeerUsername[:], buf[33:65])
    copy(dg.PeerServerAddress[:], buf[65:97])
    dg.Command = buf[97]
    copy(dg.Arguments[:], buf[98:354])
    copy(dg.Counter[:], buf[354:358])
}

// handleConnection reads datagrams from the connection and sends them to the SessionManager
func (m *SessionManager) handleConnection(conn net.Conn) {
    buf := make([]byte, 390) // Create a buffer with the size of the Datagram

    // Read the full datagram into the buffer
    _, err := io.ReadFull(conn, buf)
    if err != nil {
        if err == io.EOF {
            fmt.Println("Connection closed by client")
            return // No need to close the connection here
        }
        fmt.Printf("Error reading datagram: %v\n", err)
        conn.Close() // Close the connection in case of other errors
        return
    }

    clientOrServer := buf[0] // Read the ClientOrServer byte from the buffer

    // Step 1: Authenticate and decrypt the datagram
    if err := authenticateAndDecrypt(&buf); err != nil {
        fmt.Printf("Authentication and decryption failed: %v\n", err)

        // If it is a client session, send a generic error response
        if clientOrServer == 0 { // Check if it's a client session
            // Inline sending the error response
            if _, writeErr := conn.Write([]byte{1}); writeErr != nil {
                fmt.Printf("Failed to send error response: %v\n", writeErr)
            }
        }

        conn.Close() // Close the connection if authentication fails
        return
    }

    if clientOrServer == 0 { // Client session
        // Create and populate a ClientSession
        clientSession := &ClientSession{
            Conn: conn,
        }
        bytesToDatagram(&clientSession.Datagram, buf) // Populate using the modified buffer
        m.sessionCh <- clientSession
        // Connection remains open for client sessions
    } else { // Server session
        // Create and populate a ServerSession
        serverSession := &ServerSession{}
        bytesToDatagram(&serverSession.Datagram, buf) // Populate using the modified buffer
        m.sessionCh <- serverSession
        conn.Close() // Close the connection for server sessions
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
