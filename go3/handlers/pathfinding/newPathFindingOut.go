package pathfinding

import (
    "log"
    "time"
    "ripple/config"
    "ripple/main"
)

// NewPathFindingOut handles the command to initiate a new pathfinding request.
func NewPathFindingOut(session main.Session) {
    datagram := session.Datagram

    // Validate the counter using the ValidateCounter function from pathfinding package
    if err := ValidateCounter(datagram); err != nil {
        log.Printf("Counter validation failed for user %s: %v", datagram.Username, err)
        return // Simply return if the counter is invalid; no response is sent
    }
