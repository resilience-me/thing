package payments

import (
    "crypto/sha256"
    "fmt"
)

// GeneratePaymentOutIdentifier constructs and hashes the identifier for an outgoing payment.
func GeneratePaymentOutIdentifier(dg *Datagram) []byte {
    var buf []byte

    // For outgoing payments, the buyer is generating the identifier
    buf = append(buf, PadTo32Bytes(dg.Username)...)
    buf = append(buf, PadTo32Bytes(GetServerAddress())...)
    buf = append(buf, PadTo32Bytes(dg.PeerUsername)...)
    buf = append(buf, PadTo32Bytes(dg.PeerServerAddress)...)

    // Append amount and nonce from Arguments
    buf = append(buf, dg.Arguments[0:8]...)

    // Compute SHA-256 hash of the concatenated byte slice
    hash := sha256.Sum256(buf)

    // Return the hash as a byte slice
    return hash[:]
}

// GeneratePaymentInIdentifier constructs and hashes the identifier for an incoming payment.
func GeneratePaymentInIdentifier(dg *Datagram) []byte {
    var buf []byte

    // For incoming payments, the seller is generating the identifier
    buf = append(buf, PadTo32Bytes(dg.PeerUsername)...)
    buf = append(buf, PadTo32Bytes(dg.PeerServerAddress)...)
    buf = append(buf, PadTo32Bytes(dg.Username)...)
    buf = append(buf, PadTo32Bytes(GetServerAddress())...)

    // Append amount and nonce from Arguments
    buf = append(buf, dg.Arguments[0:8]...)

    // Compute SHA-256 hash of the concatenated byte slice
    hash := sha256.Sum256(buf)

    // Return the hash as a byte slice
    return hash[:]
}
