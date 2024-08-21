package main

import (
	"fmt"
	"sync"
)

type Conn struct {
	*net.UDPConn         // Embedding UDPConn to promote its methods
	addr *net.UDPAddr    // Keeping addr as a separate field
}

type Session struct {
	Datagram *Datagram // The datagram associated with this session
	Conn     *Conn     // Pointer to Conn; can be nil if not applicable
}

// SessionManager manages sessions and their state
type SessionManager struct {
	activeHandlers map[string]bool
	queues         map[string][]*Session
	mu             sync.Mutex
}

// NewSessionManager creates a new SessionManager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		activeHandlers: make(map[string]bool),
		queues:         make(map[string][]*Session),
	}
}

// RouteSession routes a new session or queues it if a handler is already active
func (sm *SessionManager) RouteSession(session *Session) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	username := session.Datagram.Username
	if !sm.activeHandlers[username] {
		// No active handler, process session immediately
		sm.activeHandlers[username] = true
		go sm.handleSession(session)
	} else {
		// Active handler exists, queue the session
		sm.queues[username] = append(sm.queues[username], session)
	}
}

// CloseSession processes the next session in the queue after a session finishes
func (sm *SessionManager) CloseSession(username string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if queue, exists := sm.queues[username]; exists && len(queue) > 0 {
		// Process the next session in the queue
		nextSession := queue[0]
		sm.queues[username] = queue[1:]
		go sm.handleSession(nextSession)
	} else {
		// No more sessions in the queue, mark handler as inactive
		delete(sm.activeHandlers, username)
	}
}

// handleSession processes a session and then triggers the next one
func (sm *SessionManager) handleSession(session *Session) {
	datagram := session.Datagram
	defer sm.CloseSession(datagramm.Username)

	// Handle the session here (processing logic)
	fmt.Printf("Handling session for user: %s\n", datagram.Username)

	// If this is a client connection and in the range 0x00 to 0x3F, check that peer account exists
	if datagram.Command&0xC0 == 0 { // Bit 7 (MSB) is 0: Client connection, bit 6 is 0: validate peer account exists
	    if errorMessage, err := auth.ValidatePeerExists(datagram); err != nil {
	        log.Printf("Error vaslidating peer existence for user %s: %v", datagram.Username, err)
	        main.SendErrorResponse(errorMessage, session.Conn)
	        return
	    }
	}

	// Simulate session handling
	// For example, send a datagram, wait for an ack, etc.
}
