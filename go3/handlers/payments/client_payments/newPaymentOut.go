package client_payments

import (
    "log"
    "ripple/handlers/payments"
    "ripple/main"
)

// NewPaymentOut handles the command to initiate a new outgoing payment.
func NewPaymentOut(session main.Session) {

    // Generate the payment identifier
    paymentIdentifier := payments.GeneratePaymentOutIdentifier(dg)

    // Log the identifier (for example)
    log.Printf("Generated Payment Identifier: %x\n", paymentIdentifier)

    // Step 1: Initiate the outgoing payment using the extracted username and paymentIdentifier
    err := session.PathManager.InitiateOutgoingPayment(username, paymentIdentifier)
    if err != nil {
        // Handle the error (e.g., log it, return it to the caller)
        log.Printf("failed to initiate outgoing payment for user %s: %w", username, err)
        return
    }

}
