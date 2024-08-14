package payments

import (
    "crypto/sha256"
    "ripple/main"
)

// ConcatenateAndPadAndHash takes four strings and an 8-byte slice, pads each string to 32 bytes,
// concatenates them with the 8-byte slice, and then hashes the result using SHA-256.
func ConcatenateAndPadAndHash(s1, s2, s3, s4 string, b [8]byte) []byte {
	const stringLength = 32

	// Format and pad the strings, convert to byte slice
	paddedStrings := []byte(fmt.Sprintf(
		"%-32s%-32s%-32s%-32s",
		s1,
		s2,
		s3,
		s4,
	))

	// Append the 8-byte slice to the byte slice
	concatenated := append(paddedStrings, b[:]...)

	// Compute SHA-256 hash of the concatenated result
	hash := sha256.Sum256(concatenated)

	return hash[:]
}

// ConcatenateAndPad takes four strings and one byte slice, pads each to 32 bytes, and concatenates them.
func ConcatenateAndPad(s1, s2, s3, s4 string, b []byte) string {
	const length = 32

	// Pad each string and convert the byte slice to a string with appropriate padding.
	return fmt.Sprintf(
		"%-32s%-32s%-32s%-32s%-32s",
		s1,
		s2,
		s3,
		s4,
		string(b[:min(len(b), length)])) // Ensure byte slice is not larger than 32 bytes
}



// PadUserIdentifiers pads and returns the four components needed for identifier generation.
func PadUserIdentifiers(dg *Datagram) ([]byte, []byte, []byte, []byte) {
    username := main.PadTo32Bytes(dg.Username)
    serverAddress := main.PadTo32Bytes(GetServerAddress())
    peerUsername := main.PadTo32Bytes(dg.PeerUsername)
    peerServerAddress := main.PadTo32Bytes(dg.PeerServerAddress)
    return username, serverAddress, peerUsername, peerServerAddress
}

// generatePaymentIdentifier uses nested append calls to concatenate userX, userY, and Arguments before hashing.
func generatePaymentIdentifier(userX, userY []byte, arguments []byte) []byte {
    // Concatenate userX, userY, and arguments[0:8] using nested append
    preimage := append(append(userX, userY...), arguments[0:8]...)

    // Compute SHA-256 hash of the combined byte slice
    hash := sha256.Sum256(preimage)

    // Return the hash as a byte slice
    return hash[:]
}

// Wrapper functions for outgoing and incoming payments
func GeneratePaymentOutIdentifier(dg *Datagram) []byte {
    username, serverAddress, peerUsername, peerServerAddress := PadUserIdentifiers(dg)
    userX := append(username, serverAddress...)
    userY := append(peerUsername, peerServerAddress...)
    return generatePaymentIdentifier(userX, userY, dg.Arguments)
}

func GeneratePaymentInIdentifier(dg *Datagram) []byte {
    username, serverAddress, peerUsername, peerServerAddress := PadUserIdentifiers(dg)
    userX := append(peerUsername, peerServerAddress...)
    userY := append(username, serverAddress...)
    return generatePaymentIdentifier(userX, userY, dg.Arguments)
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
