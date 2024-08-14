package client_payments

import (
    "ripple/main"
    "ripple/pathfinding"
)

type PaymentDetails struct {
    Counterpart PeerAccount
    InOrOut     byte
    Amount      uint32
    Nonce       uint32
}

// SerializePaymentDetails constructs a byte array directly from PaymentDetails.
func SerializePaymentDetails(details *PaymentDetails) []byte {
    // Create a buffer with the exact size needed
    buffer := make([]byte, 32+32+1+4+4)  // Total size: 73 bytes

    // Copy Username into buffer, ensuring it's exactly 32 bytes
    copy(buffer[0:32], details.Counterpart.Username)

    // Copy ServerAddress into buffer, ensuring it's exactly 32 bytes
    copy(buffer[32:64], details.Counterpart.ServerAddress)

    // Set InOrOut byte
    buffer[64] = details.InOrOut

    // Serialize Amount and Nonce with correct endianess directly into the buffer
    binary.BigEndian.PutUint32(buffer[65:69], details.Amount)
    binary.BigEndian.PutUint32(buffer[69:73], details.Nonce)

    return buffer
}

// GetPaymentDetails fetches the payment details from the account, including the related Path.
func GetPaymentDetails(session main.Session) *PaymentDetails {
    // Retrieve the username from the session datagram
    username := session.Datagram.Username

    // Use the existing Find method from PathManager to retrieve the account
    account := session.PathManager.Find(username)
    if account == nil || account.Payment == nil {
        return nil // Return nil if no account or no payment is found
    }

    // Find the Path using the identifier in the Payment
    path := account.Find(account.Payment.Identifier)
    if path == nil {
        return nil // Return nil if no Path is found for the payment
    }

    // Construct the PaymentDetails struct with the necessary information
    paymentDetails := &PaymentDetails{
        Counterpart: account.Payment.Counterpart,
        InOrOut:     account.Payment.InOrOut,
        Amount:      path.Amount,
        Nonce:       path.CounterIn,
    }

    return paymentDetails
}
