package client_payments

import (
    "log"
    "ripple/handlers/payments"
    "ripple/types"
    "ripple/comm"

)

// GetPayment handles the command to retrieve payment parameters.
func GetPayment(session types.Session) {

    // Extract username from the datagram
    username := session.Datagram.Username

    // Retrieve and serialize payment details using the wrapper method
    paymentDetails := payments.FetchAndSerializePaymentDetails(username)
    if paymentDetails == nil {
        paymentDetails = []byte{}  // Send an empty response if no payment details
    }

    // Send the payment details as a success response
    if err := comm.SendSuccessResponse(session.Addr, paymentDetails); err != nil {
        log.Printf("Failed to send payment details to client for user %s: %v", username, err)
        return
    }

    log.Printf("Sent payment details successfully to client for user %s.", username)
}
