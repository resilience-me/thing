package server_payments

import "ripple/payments/payment_operations"

// FindPathOut processes a pathfinding request from the buyer to the seller
func FindPathOut(session *Session) {
    payment_operations.FindPath(session.Datagram, types.Outgoing)
}
