package payments

import (
	"crypto/sha256"
	"fmt"
	"ripple/main"
)

// ConcatenateAndPadAndHash takes four strings and an 8-byte slice, pads each string to 32 bytes,
// concatenates them with the 8-byte slice, and then hashes the result using SHA-256.
func ConcatenateAndPadAndHash(s1, s2, s3, s4 string, b []byte) []byte {
	// Create a buffer of size 136 bytes to hold the four 32-byte padded strings plus 8 bytes
	buffer := make([]byte, 136)

	// Copy each string into the buffer; any overflow will be safely overwritten by the next string
	copy(buffer[0:], s1)
	copy(buffer[32:], s2)
	copy(buffer[64:], s3)
	copy(buffer[96:], s4)

	// Copy the 8-byte slice into the buffer at position 128
	copy(buffer[128:], b)

	// Compute SHA-256 hash of the buffer
	hash := sha256.Sum256(buffer)

	return hash[:]
}

// GeneratePaymentOutIdentifier generates a payment identifier for outgoing payments and returns it as a hexadecimal string.
func GeneratePaymentOutIdentifier(dg *Datagram) string {
	hash := ConcatenateAndPadAndHash(dg.Username, GetServerAddress(), dg.PeerUsername, dg.PeerServerAddress, dg.Arguments[:8])
	return fmt.Sprintf("%x", hash)
}

// GeneratePaymentInIdentifier generates a payment identifier for incoming payments and returns it as a hexadecimal string.
func GeneratePaymentInIdentifier(dg *Datagram) string {
	hash := ConcatenateAndPadAndHash(dg.PeerUsername, dg.PeerServerAddress, dg.Username, GetServerAddress(), dg.Arguments[:8])
	return fmt.Sprintf("%x", hash)
}

// GeneratePayment creates a Payment struct based on the provided datagram and inOrOut value.
func generatePayment(datagram *Datagram, inOrOut byte) *Payment {
    // Create the counterpart PeerAccount using the two relevant fields from the datagram
    counterpart := PeerAccount{
        Username:      datagram.PeerUsername,
        ServerAddress: datagram.PeerServerAddress,
    }

    // Generate the payment identifier
    paymentIdentifier := GeneratePaymentOutIdentifier(datagram)

    // Initialize and return the Payment struct
    return &Payment{
        Identifier:  paymentIdentifier,
        Counterpart: counterpart,
        InOrOut:     inOrOut,
    }
}

// GenerateOutgoingPayment generates a Payment struct for an outgoing payment.
func generateOutgoingPayment(datagram *Datagram) *Payment {
    return GeneratePayment(datagram, 0) // 0 for outgoing
}

// GenerateIncomingPayment generates a Payment struct for an incoming payment.
func generateIncomingPayment(datagram *Datagram) *Payment {
    return GeneratePayment(datagram, 1) // 1 for incoming
}

// GenerateAndInitiatePaymentOut handles the generation of the payment identifier and initiation of the outgoing payment.
func GenerateAndInitiatePaymentOut(session main.Session) error {
    // Generate the Payment struct for an outgoing payment
    payment := GenerateOutgoingPayment(session.Datagram)

    // Initiate the outgoing payment using the constructed Payment struct
    session.PathManager.initiatePayment(session.Datagram.Username, payment)

    return nil
}

// GenerateAndInitiatePaymentIn handles the generation of the payment identifier and initiation of the incoming payment.
func GenerateAndInitiatePaymentIn(session main.Session) error {
    // Generate the Payment struct for an incoming payment
    payment := GenerateIncomingPayment(session.Datagram)

    // Initiate the incoming payment using the constructed Payment struct
    session.PathManager.initiatePayment(session.Datagram.Username, payment)

    return nil
}

func (account *Account) CalculateCommittedAmounts() map[string]CommitTotals {
    now := time.Now().Unix()
    commitExpiration := now - int64(config.CommitTimeout.Seconds()) // Pre-calculate the expiration threshold

    totals := make(map[string]CommitTotals)
    
    for _, path := range account.Paths {
        if path.Commit && path.Timestamp.Unix() > commitExpiration { // Check if the path is still valid
            // Update totals for outgoing transactions (credit given by this account)
            outgoingPeer := path.Outgoing.Username
            totals[outgoingPeer].Outgoing += path.Amount

            // Update totals for incoming transactions (credit received by this account)
            incomingPeer := path.Incoming.Username
            totals[incomingPeer].Incoming += path.Amount
        }
    }
    
    return totals
}
