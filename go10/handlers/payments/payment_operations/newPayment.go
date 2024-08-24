package payment_operations

import (
    "log"                 // For logging errors and success messages
    "ripple/comm"         // For sending error and success responses to the client
    "ripple/handlers/payments"  // For calling the GenerateAndInitiatePayment function
    "ripple/types"
)

// NewPayment is a shared function to handle the payment initialization process.
func NewPayment(session types.Session, inOrOut byte) {
    // Retrieve the Datagram from the session
    datagram := session.Datagram

    // Extract username from the datagram
    username := datagram.Username

    // Generate the payment identifier and initiate the payment
    err := payments.GenerateAndInitiatePayment(datagram, inOrOut)
    if err != nil {
        log.Printf("Failed to initiate payment for user %s: %v", username, err)
        comm.SendErrorResponse(session.Addr, "Failed to initiate payment.")
        return
    }
    log.Printf("Payment initialized successfully for user %s.", username)

    // Send success response
    if err := comm.SendSuccessResponse(session.Addr, []byte("Payment initialized successfully.")); err != nil {
        log.Printf("Failed to send success response to user %s: %v", username, err)
        return
    }

    log.Printf("Sent success response to client for user %s.", username)
}
