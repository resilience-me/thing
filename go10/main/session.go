package main

import (
	"log"
	"sync"
	"ripple/auth"
	"ripple/comm"
	"ripple/types"
)

// SessionManager manages sessions and their state
type SessionManager struct {
	activeHandlers map[string]bool
	queues         map[string][]*types.Session
	mu             sync.Mutex
	wg             sync.WaitGroup
}

// NewSessionManager creates a new SessionManager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		activeHandlers: make(map[string]bool),
		queues:         make(map[string][]*types.Session),
	}
}

// RouteSession routes a new session or queues it if a handler is already active
func (sm *SessionManager) RouteSession(session *types.Session) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.wg.Add(1)

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

	sm.wg.Done()

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
func (sm *SessionManager) handleSession(session *types.Session) {
	datagram := session.Datagram
	command := datagram.Command
	username := datagram.Username

	defer sm.CloseSession(username)

	// Handle the session here (processing logic)
	log.Printf("Handling session for user: %s\n", username)


	// If this is a client connection, check that peer account exists
	// But only if the command included a peer.
	if command&0x80 == 0 && datagram.PeerUsername != "" { // Bit 7 (MSB) is 0
	    if errorMessage, err := auth.ValidatePeerExists(datagram); err != nil {
	        log.Printf("Error validating peer existence for user %s: %v", username, err)
	        comm.SendErrorResponse(session.Addr, errorMessage)
	        return
	    }
	}
	
	handler := commandHandlers[command]
	if handler == nil {
		log.Printf("Unknown command: %d\n", command)
		return
	}
	
	handler(*session)
}
