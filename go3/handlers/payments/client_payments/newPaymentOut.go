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

    // Validate the counter using the ValidateCounter function from payments package
    if err := payments.ValidateCounter(datagram); err != nil {
        log.Printf("Counter validation failed for user %s: %v", username, err)
        main.SendErrorResponse("Received counter is not valid.", session.Conn)
        return
    }
    
    // Generate the payment identifier and initiate the outgoing payment
    err := payments.GenerateAndInitiatePaymentOut(session, datagram, username)
    if err != nil {
        log.Printf("Failed to initiate outgoing payment for user %s: %v", username, err)
        main.SendErrorResponse("Failed to initiate payment.", session.Conn)
        return
    }
    log.Printf("Payment initialized successfully for user %s.", username)

    // Write the new client-side counter value using the setter in db_payments
    if err := db_payments.SetCounter(datagram, datagram.Counter); err != nil {
        log.Printf("Error writing counter to file for user %s: %v", username, err)
        main.SendErrorResponse("Failed to write counter.", session.Conn)
        return
    }

    // Send success response
    if err := main.SendSuccessResponse([]byte("Payment initialized successfully."), session.Conn); err != nil {
        log.Printf("Failed to send success response to user %s: %v", username, err)
        return
    }
    log.Printf("Sent success response to client for user %s.", username)
}
