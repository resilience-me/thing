package server_payments

import (
    "ripple/types"
    "ripple/handlers/payments/payment_operations"
)

// FindPathIn processes a pathfinding request from the seller to the buyer
func FindPathIn(session types.Session) {
    payment_operations.FindPath(session.Datagram, types.Incoming)
}
