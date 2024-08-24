package server_payments

import (
    "ripple/main"
    "ripple/types"
    "ripple/handlers/payments/payment_operations"
)

// FindPathOut processes a pathfinding request from the buyer to the seller
func FindPathOut(session main.Session) {
    payment_operations.FindPath(session.Datagram, types.Outgoing)
}
