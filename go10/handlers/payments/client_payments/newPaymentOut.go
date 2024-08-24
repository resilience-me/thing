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
    
    // Generate the payment identifier and initiate the outgoing payment
    err := payments.GenerateAndInitiatePaymentOut(session, datagram, username)
    if err != nil {
        log.Printf("Failed to initiate outgoing payment for user %s: %v", username, err)
        comm.SendErrorResponse("Failed to initiate payment.", session.Addr)
        return
    }
    log.Printf("Payment initialized successfully for user %s.", username)

    // Send success response
    if err := comm.SendSuccessResponse([]byte("Payment initialized successfully."), session.Addr); err != nil {
        log.Printf("Failed to send success response to user %s: %v", username, err)
        return
    }

    log.Printf("Sent success response to client for user %s.", username)
}
