package server_payments

import "ripple/payments/payment_operations"

// FindPathIn processes a pathfinding request from the seller to the buyer
func FindPathIn(session *Session) {
    payment_operations.FindPath(session.Datagram, types.Incoming)
}
