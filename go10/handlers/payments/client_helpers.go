package payments

import (
    "ripple/types"
    "ripple/pathfinding"
)

// serializePaymentDetails constructs a byte array from the payment details
func serializePaymentDetails(payment *pathfinding.Payment, amount uint32) []byte {
    buffer := concatNameAndServer(payment.Counterpart.Username, payment.Counterpart.ServerAddress)
    buffer = append(buffer, payment.InOrOut)
    amountAndNonce := append(types.Uint32ToBytes(path.Amount), types.Uint32ToBytes(payment.Nonce)...)
    buffer = append(buffer, amountAndNonce...)
    return buffer
}

// Wrapper function to fetch and serialize payment details
func FetchAndSerializePaymentDetails(username string) []byte {
    // Use the existing Find method from PathManager to retrieve the account
    account := pathfinding.PathManager.Find(username)
    if account == nil || account.Payment == nil {
        return nil // Return nil if no account or no payment is found
    }
    // Find the Path using the identifier in the Payment
    path := account.Find(account.Payment.Identifier)
    if path == nil {
        return nil // Return nil if no Path is found for the payment
    }

    return serializePaymentDetails(account.Payment, path.Amount)
}
