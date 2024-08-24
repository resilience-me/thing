package server_payments

import (
    "ripple/types"
    "ripple/handlers/payments/payment_operations"
)

// FindPathOut processes a pathfinding request from the buyer to the seller
func FindPathOut(session types.Session) {
    payment_operations.FindPath(session.Datagram, types.Outgoing)
}
