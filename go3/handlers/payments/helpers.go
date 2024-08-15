package payments

import (
	"crypto/sha256"
	"fmt"
	"ripple/main"
	"ripple/pathfinding"
)

// concatenateAndPadAndHash takes four strings and an 8-byte slice, pads each string to 32 bytes,
// concatenates them with the 8-byte slice, and then hashes the result using SHA-256.
func concatenateAndPadAndHash(s1, s2, s3, s4 string, b []byte) []byte {
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

// generatePaymentInIdentifier generates a payment identifier for incoming payments and returns it as a hexadecimal string.
func generatePaymentInIdentifier(dg *Datagram) string {
	hash := ConcatenateAndPadAndHash(dg.PeerUsername, dg.PeerServerAddress, dg.Username, GetServerAddress(), dg.Arguments[:8])
	return fmt.Sprintf("%x", hash)
}

// generatePaymentOutIdentifier generates a payment identifier for outgoing payments and returns it as a hexadecimal string.
func generatePaymentOutIdentifier(dg *Datagram) string {
	hash := ConcatenateAndPadAndHash(dg.Username, GetServerAddress(), dg.PeerUsername, dg.PeerServerAddress, dg.Arguments[:8])
	return fmt.Sprintf("%x", hash)
}

// generatePaymentIn generates a Payment struct for an incoming payment.
func generatePaymentIn(datagram *Datagram, identifier string) *Payment {
    return pathfinding.NewPayment(datagram, identifier, 0)
}

// generatePaymentOut generates a Payment struct for an outgoing payment.
func generatePaymentOut(datagram *Datagram, identifier string) *Payment {
    return pathfinding.NewPayment(datagram, identifier, 1)
}

// GenerateAndInitiatePaymentIn handles the generation of the payment identifier and initiation of the incoming payment.
func GenerateAndInitiatePaymentIn(session main.Session) {
    // Generate the Payment struct for an incoming payment

    paymentIdentifier := generatePaymentInIdentifier(datagram)
    payment := generatePaymentIn(session.Datagram, paymentIdentifier)

    // Initiate the incoming payment using the constructed Payment struct
    session.PathManager.initiatePayment(session.Datagram.Username, payment)
}

// GenerateAndInitiatePaymentOut handles the generation of the payment identifier and initiation of the outgoing payment.
func GenerateAndInitiatePaymentOut(session main.Session) {
    // Generate the Payment struct for an outgoing payment
    paymentIdentifier := generatePaymentOutIdentifier(datagram)
    payment := generatePaymentOut(session.Datagram, paymentIdentifier

    // Initiate the outgoing payment using the constructed Payment struct
    session.PathManager.initiatePayment(session.Datagram.Username, payment)
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
