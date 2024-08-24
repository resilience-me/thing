package client_payments

import (
    "encoding/binary"
    "ripple/main"
    "ripple/pathfinding"
)

type PaymentDetails struct {
    Counterpart PeerAccount
    Amount      uint32
    InOrOut     byte
    Nonce       uint32
}

// NewPaymentDetails is a constructor for creating a NewPaymentDetails struct
func NewPaymentDetails(counterpart PeerAccount, amount uint32, inOrOut byte, nonce uint32) *PaymentDetails {
    return &PaymentDetails{
        Counterpart: counterpart,
        Amount:      amount,
        InOrOut:     inOrOut,
        Nonce:       nonce,
    }
}

// serializePaymentDetails constructs a byte array directly from PaymentDetails.
func serializePaymentDetails(details *PaymentDetails) []byte {
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

// getPaymentDetails fetches the payment details from the account, including the related Path.
func getPaymentDetails(username string) *PaymentDetails {
    // Use the existing Find method from PathManager to retrieve the account
    if account := pathfinding.PathManager.Find(username); account != nil && account.Payment != nil {
        // Find the Path using the identifier in the Payment
        if path := account.Find(account.Payment.Identifier); path != nil {
            return NewPaymentDetails(account.Payment.Counterpart, path.Amount, account.Payment.InOrOut, account.Payment.Nonce)
        }
    }
    return nil // Return nil if no account or no payment is found
}

// Wrapper function to fetch and serialize payment details
func fetchAndSerializePaymentDetails(username string) []byte {
    // Fetch the payment details
    details := getPaymentDetails(username)
    if details == nil {
        return nil  // Return nil to indicate no data is available
    }

    // Serialize the payment details
    serializedData := serializePaymentDetails(details)
    return serializedData
}
