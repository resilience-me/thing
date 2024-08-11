package main

import (
    "fmt"  // For formatted I/O operations like Printf, Println, etc.
    "io"   // For I/O operations, including ReadFull
    "net"  // For networking operations, including net.Conn, net.Listen, etc.
    "os"   // For OS-level operations like exiting the program with os.Exit
)

// Session struct represents a network session with an optional connection
type Session struct {
    Datagram  *Datagram  // Using the same name for the identifier and the type
    Conn      net.Conn   // Connection, can be nil if no network connection is associated
}

// SessionManager manages the processing of sessions
type SessionManager struct {
    sessionCh      chan Session
    closedCh       chan string
    activeHandlers map[string]bool
    queues         map[string][]Session
}

// run method for the SessionManager
func (m *SessionManager) run() {
    for {
        select {
        case session := <-m.sessionCh:
            username := session.Datagram.Username

            if !m.activeHandlers[username] {
                m.activeHandlers[username] = true
                go m.handleSession(session)
            } else {
                m.queues[username] = append(m.queues[username], session)
            }

        case username := <-m.closedCh:
            if queue, exists := m.queues[username]; exists && len(queue) > 0 {
                nextSession := queue[0]
                m.queues[username] = queue[1:]
                go m.handleSession(nextSession)
            } else {
                delete(m.activeHandlers, username)
            }
        }
    }
}

// handleSession manages the processing of a session's datagram
func (m *SessionManager) handleSession(session Session) {
    defer func() {
        if session.Conn != nil { // Check if the connection exists
            session.Conn.Close() // Close the connection to free up resources
        }
        m.closedCh <- session.Datagram.Username // Notify that the session is closed
    }()

    command := session.Datagram.Command
    handler := commandHandlers[command]
    if handler == nil {
        fmt.Printf("Unknown command: %d\n", command)
        return
    }

    handler(session)
}

// handleClientConnection processes a connection from a client.
func (m *SessionManager) handleClientConnection(conn net.Conn, buf []byte) {
    dg, errorMessage, err := validateAndParseClientDatagram(buf)
    if err != nil {
        conn.Write([]byte{1})                      // Indicate error with '1'
        conn.Write([]byte(errorMessage))           // Send the specific error message
        fmt.Printf("Error during datagram validation: %v\n", err)
        conn.Close()
        return
    }
    // Send the session to the session channel for further processing
    m.sessionCh <- Session{Datagram: dg, Conn: conn} // Client sessions keep the connection open
}

// handleServerConnection processes a connection from another server.
func (m *SessionManager) handleServerConnection(buf []byte) {
    dg, err := validateAndParseServerDatagram(buf)
    if err != nil {
        fmt.Printf("Error validating server datagram: %v\n", err)
        return
    }

    // Send the session to the session channel for further processing
    m.sessionCh <- Session{Datagram: dg} // Conn is nil for server sessions
}

// handleConnection reads datagrams from the connection and decides whether to handle a client or server connection.
func (m *SessionManager) handleConnection(conn net.Conn) {
    buf := make([]byte, 389) // Adjust the buffer size according to your actual data size
    _, err := io.ReadFull(conn, buf)
    if err != nil {
        fmt.Printf("Error reading datagram: %v\n", err)
        conn.Close()
        return
    }

    // Determine whether it's a client or server session and handle accordingly
    if buf[0]&0x80 == 0 { // Client session if MSB is 0
        m.handleClientConnection(conn, buf)
    } else { // Server session
        m.handleServerConnection(buf)
        conn.Close() // Close the connection directly after processing
    }
}

// Main function with inlined server logic
func main() {
    manager := &SessionManager{
        sessionCh:      make(chan Session),
        closedCh:       make(chan string),
        activeHandlers: make(map[string]bool),
        queues:         make(map[string][]Session),
    }
    go manager.run()

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
        go manager.handleConnection(conn)
    }
}
