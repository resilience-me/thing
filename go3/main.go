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

// handleConnection reads datagrams from the connection and sends them to the SessionManager
func (m *SessionManager) handleConnection(conn net.Conn) {
    buf := make([]byte, 389) // Adjust the buffer size according to your actual data size
    _, err := io.ReadFull(conn, buf)
    if err != nil {
        fmt.Printf("Error reading datagram: %v\n", err)
        conn.Close()
        return
    }

    dg := parseDatagram(buf)
    isClientSession := dg.Command & 0x80 == 0

    if errorCode, err := CheckUserAndPeerExist(dg); err != nil {
        if errorCode != 0 && isClientSession {
            _, writeErr := conn.Write([]byte{errorCode})
            if writeErr != nil {
                fmt.Printf("Failed to send error code to client: %v\n", writeErr)
            }
        }
        fmt.Printf("Error during user and peer check: %v\n", err)
        conn.Close()
        return
    }

    // Authenticate and parse the datagram
    dg, err := validateAndParseDatagram(buf)
    if err != nil {
        fmt.Printf("Error processing incoming datagram: %v\n", err)
        conn.Close()
        return
    }

    // Prepare the session struct
    session := Session{Datagram: dg} // Conn is nil by default

    // Determine whether it's a client or server session based on the most significant bit of dg.Command
    if isClientSession { // Check if the most significant bit is 0 (client)
        session.Conn = conn // Maintain the connection open for client sessions
    } else { // Most significant bit is 1 (server)
        conn.Close() // Close the connection for server sessions
        // No need to set session.Conn to nil; it is already nil by default
    }

    // Send the session to the session channel for further processing
    m.sessionCh <- session
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
