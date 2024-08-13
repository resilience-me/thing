package client_payments

import (
    "log"
    "ripple/handlers/payments"
    "ripple/main"
)

// NewPaymentOut handles the command to initiate a new outgoing payment.
func NewPaymentOut(session main.Session) {

    // Retrieve the Datagram from the session
    datagram := session.Datagram

    // Extract username from the datagram
    username := datagram.Username
    
    // Generate the payment identifier
    paymentIdentifier := payments.GeneratePaymentOutIdentifier(datagram)

    // Log the identifier (for example)
    log.Printf("Generated Payment Identifier: %x\n", paymentIdentifier)

    // Step 1: Initiate the outgoing payment using the extracted username and paymentIdentifier
    err := session.PathManager.InitiateOutgoingPayment(username, paymentIdentifier)
    if err != nil {
        // Handle the error (e.g., log it, return it to the caller)
        log.Printf("Failed to initiate outgoing payment for user %s: %v", username, err)
        return
    }
    log.Printf("Payment initialized successfully for user %s.", username)

    // Send success response
    if err := main.SendSuccessResponse([]byte("Payment initialized successfully."), session.Conn); err != nil {
        log.Printf("Failed to send success response to user %s: %v", username, err)
        return
    }
    log.Printf("Sent success response to client for user %s.", username)
}
