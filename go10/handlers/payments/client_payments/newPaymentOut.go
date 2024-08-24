package client_payments

import (
    "ripple/types"
    "ripple/handlers/payments/payment_operations"
)

// NewPaymentOut handles the command to initiate a new outgoing payment.
func NewPaymentOut(session types.Session) {
    payment_operations.NewPayment(session, types.Outgoing)
}
