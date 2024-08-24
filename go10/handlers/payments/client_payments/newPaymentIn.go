package client_payments

import (
    "ripple/types"
    "ripple/handlers/payments/payment_operations"
)

// NewPaymentIn handles the command to initiate a new incoming payment.
func NewPaymentIn(session types.Session) {
    payment_operations.NewPayment(session, types.Incoming)
}
