package payments

import (
	"crypto/sha256"
	"fmt"
	"ripple/main"
)

// ConcatenateAndPadAndHash takes four strings and an 8-byte slice, pads each string to 32 bytes,
// concatenates them with the 8-byte slice, and then hashes the result using SHA-256.
func ConcatenateAndPadAndHash(s1, s2, s3, s4 string, b []byte) []byte {
	// Convert padded strings to byte slice
	paddedStrings := []byte(fmt.Sprintf(
		"%-32s%-32s%-32s%-32s",
		s1,
		s2,
		s3,
		s4,
	))

	// Append the byte slice to the byte slice
	concatenated := append(paddedStrings, b...)

	// Compute SHA-256 hash of the concatenated result
	hash := sha256.Sum256(concatenated)

	return hash[:]
}

// GeneratePaymentOutIdentifier generates a payment identifier for outgoing payments.
func GeneratePaymentOutIdentifier(dg *Datagram) []byte {
    // Directly pass the strings and the byte slice to ConcatenateAndPadAndHash
    return ConcatenateAndPadAndHash(dg.Username, GetServerAddress(), dg.PeerUsername, dg.PeerServerAddress, dg.Arguments[:8])
}

// GeneratePaymentInIdentifier generates a payment identifier for incoming payments.
func GeneratePaymentInIdentifier(dg *Datagram) []byte {
    // Directly pass the strings and the byte slice to ConcatenateAndPadAndHash
    return ConcatenateAndPadAndHash(dg.PeerUsername, dg.PeerServerAddress, dg.Username, GetServerAddress(), dg.Arguments[:8])
}


// GenerateAndInitiatePaymentOut handles the generation of the payment identifier and initiation of the outgoing payment.
func GenerateAndInitiatePaymentOut(session main.Session, datagram *Datagram, username string) error {

    // Generate the payment identifier
    paymentIdentifier := GeneratePaymentOutIdentifier(datagram)

    // Initiate the outgoing payment using the extracted username and paymentIdentifier
    err := session.PathManager.InitiateOutgoingPayment(username, paymentIdentifier)
    if err != nil {
        return fmt.Errorf("Failed to initiate outgoing payment for user %s: %v", username, err)
    }

    return nil
}
