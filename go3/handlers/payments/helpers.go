package payments

import (
    "crypto/sha256"
    "ripple/main"
)

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
    userX := append(main.PadTo32Bytes(dg.Username), main.PadTo32Bytes(GetServerAddress())...)
    userY := append(main.PadTo32Bytes(dg.PeerUsername), main.PadTo32Bytes(dg.PeerServerAddress)...)
    return generatePaymentIdentifier(userX, userY, dg.Arguments)
}

func GeneratePaymentInIdentifier(dg *Datagram) []byte {
    userX := append(main.PadTo32Bytes(dg.PeerUsername), main.PadTo32Bytes(dg.PeerServerAddress)...)
    userY := append(main.PadTo32Bytes(dg.Username), main.PadTo32Bytes(GetServerAddress())...)
    return generatePaymentIdentifier(userX, userY, dg.Arguments)
}

