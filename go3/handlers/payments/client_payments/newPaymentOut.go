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
    
    // Generate the payment identifier
    paymentIdentifier := payments.GeneratePaymentOutIdentifier(datagram)

    // Log the identifier (for example)
    log.Printf("Generated Payment Identifier: %x\n", paymentIdentifier)

    // Initiate the outgoing payment using the extracted username and paymentIdentifier
    err := session.PathManager.InitiateOutgoingPayment(username, paymentIdentifier)
    if err != nil {
        // Handle the error (e.g., log it, return it to the caller)
        log.Printf("Failed to initiate outgoing payment for user %s: %v", username, err)
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
