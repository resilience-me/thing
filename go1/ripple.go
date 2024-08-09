package main

import (
    "fmt"
    "io"
    "net"
    "os"
)

// BaseSession struct, embedding the Datagram
type BaseSession struct {
    Datagram
}

// Session interface with a GetDatagram method
type Session interface {
    GetDatagram() *Datagram
}

// GetDatagram method for BaseSession
func (bs *BaseSession) GetDatagram() *Datagram {
    return &bs.Datagram
}

// ClientSession struct
type ClientSession struct {
    BaseSession
    Conn net.Conn
}

// ServerSession struct
type ServerSession struct {
    BaseSession
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

// bytesToDatagram populates a Datagram struct from a byte slice
func bytesToDatagram(dg *Datagram, buf []byte) {
    dg.ClientOrServer = buf[0]
    dg.Username = ToString(buf[1:33])
    dg.PeerUsername = ToString(buf[33:65])
    dg.PeerServerAddress = ToString(buf[65:97])
    dg.Command = buf[97]
    copy(dg.Arguments[:], buf[98:354])
    copy(dg.Counter[:], buf[354:358])
}

// handleConnection reads datagrams from the connection and sends them to the SessionManager
func (m *SessionManager) handleConnection(conn net.Conn) {
    buf := make([]byte, 390)

    _, err := io.ReadFull(conn, buf)
    if err != nil {
        if err == io.EOF {
            fmt.Println("Connection closed by client")
            return
        }
        fmt.Printf("Error reading datagram: %v\n", err)
        conn.Close()
        return
    }

    clientOrServer := buf[0]
    if err := authenticateAndDecrypt(&buf); err != nil {
        fmt.Printf("Authentication and decryption failed: %v\n", err)
        if clientOrServer == 0 {
            if _, writeErr := conn.Write([]byte{1}); writeErr != nil {
                fmt.Printf("Failed to send error response: %v\n", writeErr)
            }
        }
        conn.Close()
        return
    }

    if clientOrServer == 0 {
        clientSession := &ClientSession{Conn: conn}
        bytesToDatagram(&clientSession.Datagram, buf)
        m.sessionCh <- clientSession
    } else {
        serverSession := &ServerSession{}
        bytesToDatagram(&serverSession.Datagram, buf)
        m.sessionCh <- serverSession
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
