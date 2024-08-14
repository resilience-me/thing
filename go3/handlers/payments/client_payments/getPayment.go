package client_payments

import (
    "log"
    "ripple/handlers/payments"
    "ripple/main"
)

// GetPayment handles the command to retrieve payment parameters.
func GetPayment(session main.Session) {
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

    // Retrieve and serialize payment details using the wrapper method
    paymentDetails := fetchAndSerializePaymentDetails(session)
    if paymentDetails == nil {
        paymentDetails = []byte{}  // Send an empty response if no payment details
    }

    // Send the payment details as a success response
    if err := main.SendSuccessResponse(paymentDetails, session.Conn); err != nil {
        log.Printf("Failed to send payment details to client for user %s: %v", username, err)
        return
    }

    log.Printf("Sent payment details successfully to client for user %s.", username)
}
