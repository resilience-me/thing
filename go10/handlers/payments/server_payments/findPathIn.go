package server_payments

import (
    "ripple/main"
    "ripple/types"
    "ripple/handlers/payments/payment_operations"
)

// FindPathIn processes a pathfinding request from the seller to the buyer
func FindPathIn(session main.Session) {
    payment_operations.FindPath(session.Datagram, types.Incoming)
}
