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
	shutdown       chan struct{}
	wg             sync.WaitGroup // WaitGroup to track active handlers and queued sessions
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
	    m.wg.Done()
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

// handleConnection reads datagrams from the connection and decides whether to handle a client or server connection.
func (m *SessionManager) handleConnection(conn net.Conn) {
    buf := make([]byte, 389) // Adjust the buffer size according to your actual data size
    _, err := io.ReadFull(conn, buf)
    if err != nil {
        fmt.Printf("Error reading datagram: %v\n", err)
        conn.Close()
	m.wg.Done()
        return
    }

    dg := parseDatagram(buf)
    session := Session{Datagram: dg}
    
    // Determine whether it's a client or server session
    if dg.Command & 0x80 == 0 { // Client session if MSB is 0
        errorMessage, err := validateClientDatagram(buf, dg)
        if err != nil {
            SendErrorResponse(errorMessage, conn)
            fmt.Printf("Error during datagram validation: %v\n", err)
            conn.Close()
	    m.wg.Done()
            return
        }
        session.Conn = conn // Prepare session with connection for clients
    } else { // Server session
        conn.Close() // Close the connection directly
        if err := validateServerDatagram(buf, dg); err != nil {
            fmt.Printf("Error validating server datagram: %v\n", err)
	    m.wg.Done()
            return
        }
    }

    // Send the session to the session channel for further processing
    m.sessionCh <- session
}

func (m *SessionManager) shutdownHandler(listener net.Listener) {
	interruptCount := 0 // Scoped to this function

	// Channel to capture OS signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for range signals {
		interruptCount++
		if interruptCount == 1 {
			fmt.Println("Interrupt received, initiating graceful shutdown...")
			close(manager.shutdown)  // Signal to shutdown the manager and other components
			listener.Close()         // Close the listener to stop accepting new connections
			continue // Skip to the next iteration
		}

		if interruptCount >= 9 {
			fmt.Println("Force quitting after 9 interrupts...")
			os.Exit(1) // Force exit after receiving 9 interrupts
		}

		fmt.Printf("Interrupt received (%d/9), press Ctrl+C again to force quit...\n", interruptCount)
	}
}

// Main function with inlined server logic
func main() {
    manager := &SessionManager{
        sessionCh:      make(chan Session),
        closedCh:       make(chan string),
        activeHandlers: make(map[string]bool),
        queues:         make(map[string][]Session),
        shutdown:       make(chan struct{}),
    }
    go manager.run()

    listener, err := net.Listen("tcp", ":2012")
    if err != nil {
        fmt.Printf("Error starting TCP server: %v\n", err)
        os.Exit(1)
    }

    // Goroutine to handle shutdown
    go manager.shutdownHandler(listener)

    fmt.Println("Listening on port 2012...")
    // Loop to handle connections with a select for shutdown
    for {
        select {
        case <-manager.shutdown:
            fmt.Println("Server is shutting down...")
            return // Exit the main loop and function

        default:
            conn, err := listener.Accept()
            if err != nil {
                fmt.Printf("Error accepting connection: %v\n", err)
                continue // Directly continue in case of error
            }
	    manager.wg.Add(1)
            go manager.handleConnection(conn)
        }
    }

    // Wait for all sessions to finish before exiting
    manager.wg.Wait()
    fmt.Println("All sessions and queues have been processed. Exiting.")
}
