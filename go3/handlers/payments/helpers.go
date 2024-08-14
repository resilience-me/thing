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
