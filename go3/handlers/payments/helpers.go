package payments

import (
    "crypto/sha256"
    "ripple/main"
)

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
func GenerateAndInitiatePaymentOut(session main.Session, datagram *Datagram, username string) ([]byte, error) {

    // Generate the payment identifier
    paymentIdentifier := GeneratePaymentOutIdentifier(datagram)

    // Log the identifier
    log.Printf("Generated Payment Identifier: %x\n", paymentIdentifier)

    // Initiate the outgoing payment using the extracted username and paymentIdentifier
    err := session.PathManager.InitiateOutgoingPayment(username, paymentIdentifier)
    if err != nil {
        log.Printf("Failed to initiate outgoing payment for user %s: %v", username, err)
        return nil, err
    }

    log.Printf("Payment initialized successfully for user %s.", username)
    return paymentIdentifier, nil
}
