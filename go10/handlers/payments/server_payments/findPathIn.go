package server_payments

import "ripple/payment_operations"

// FindPathIn processes a pathfinding request from the seller to the buyer
func FindPathIn(session *Session) {
    payment_operations.HandleFindPath(session.Datagram, types.Incoming)
}
