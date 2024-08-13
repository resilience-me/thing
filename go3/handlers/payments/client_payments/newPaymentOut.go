package client_payments

import (
    "fmt"
    "ripple/main"
)

// NewPaymentOut handles the command to initiate a new outgoing payment.
func NewPaymentOut(session main.Session) {

    // Step 1: Initiate the outgoing payment using the extracted username and paymentID
    err := session.PathManager.InitiateOutgoingPayment(username, paymentID)
    if err != nil {
        // Handle the error (e.g., log it, return it to the caller)
        return fmt.Errorf("failed to initiate outgoing payment for user %s: %w", username, err)
    }

}
