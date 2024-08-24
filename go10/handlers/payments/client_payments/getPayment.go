package client_payments

import (
    "log"
    "ripple/handlers/payments"
    "ripple/main"
    "ripple/database/db_client"
)

// GetPayment handles the command to retrieve payment parameters.
func GetPayment(session main.Session) {
    // Retrieve the Datagram from the session
    datagram := session.Datagram

    // Extract username from the datagram
    username := datagram.Username


    // Retrieve and serialize payment details using the wrapper method
    paymentDetails := fetchAndSerializePaymentDetails(session)
    if paymentDetails == nil {
        paymentDetails = []byte{}  // Send an empty response if no payment details
    }

    // Send the payment details as a success response
    if err := comm.SendSuccessResponse(paymentDetails, session.Addr); err != nil {
        log.Printf("Failed to send payment details to client for user %s: %v", datagram.Username, err)
        return
    }

    log.Printf("Sent payment details successfully to client for user %s.", datagram.Username)
}
