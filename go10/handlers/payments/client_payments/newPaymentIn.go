package client_payments

import (
    "ripple/main"
    "ripple/types"
    "ripple/handlers/payments/payment_operations"
)

// NewPaymentIn handles the command to initiate a new incoming payment.
func NewPaymentIn(session main.Session) {
    payment_operations.NewPayment(session, types.Incoming)
}
