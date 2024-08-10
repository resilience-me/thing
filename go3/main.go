package main

import (
    "fmt"
    "io"
    "net"
    "os"
)

// Session interface with a GetDatagram method
type Session interface {
    GetDatagram() *Datagram
}

// ClientSession struct
type ClientSession struct {
    Datagram
    Conn net.Conn
}

// ServerSession struct
type ServerSession struct {
    Datagram
}

// GetDatagram method for ClientSession
func (cs *ClientSession) GetDatagram() *Datagram {
    return &cs.Datagram
}

// GetDatagram method for ServerSession
func (ss *ServerSession) GetDatagram() *Datagram {
    return &ss.Datagram
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
            username := session.GetDatagram().Username

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
        m.closedCh <- session.GetDatagram().Username
    }()

    command := session.GetDatagram().Command
    handler := commandHandlers[command]
    if handler == nil {
        fmt.Printf("Unknown command: %d\n", command)
        return
    }

    handler(session)
}

// handleConnection reads datagrams from the connection and sends them to the SessionManager
func (m *SessionManager) handleConnection(conn net.Conn) {
    buf := make([]byte, 402) // Adjust the buffer size according to your actual data size (e.g., identifier + salt + ciphertext)
    _, err := io.ReadFull(conn, buf)
    if err != nil {
        fmt.Printf("Error reading datagram: %v\n", err)
        conn.Close()
        return
    }

    // Authenticate and parse the datagram
    dg, err := authenticateAndParseDatagram(buf)
    if err != nil {
        fmt.Printf("Error processing incoming datagram: %v\n", err)
        conn.Close()
        return
    }

    // Determine whether it's a client or server session based on the most significant bit of dg.Command
    if dg.Command & 0x80 == 0 { // Check if the most significant bit is 0 (client)
        m.sessionCh <- &ClientSession{tx, conn}
    } else { // Most significant bit is 1 (server)
        m.sessionCh <- &ServerSession{tx}
        conn.Close()
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
