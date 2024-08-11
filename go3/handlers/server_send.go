package handlers

import (
    "fmt"            // For formatted I/O
    "ripple/main"    // For the main package, which includes the Session and Datagram types
)

// CreateSignedDatagram creates a signed datagram by serializing it and adding a signature.
// It requires the session to load the secret key for HMAC generation.
func CreateSignedDatagram(session main.Session, dg *main.Datagram) ([]byte, error) {
    // Serialize the datagram without the signature field
    serializedData, err := main.SerializeDatagram(dg)
    if err != nil {
        return nil, fmt.Errorf("failed to serialize datagram: %w", err)
    }

    // Load the secret key for HMAC generation
    secretKey, err := main.LoadServerSecretKey(session.Datagram)
    if err != nil {
        return nil, fmt.Errorf("failed to load server secret key: %w", err)
    }

    // Generate HMAC for the serialized data
    signature, err := main.GenerateHMAC(serializedData, secretKey)
    if err != nil {
        return nil, fmt.Errorf("failed to generate HMAC: %w", err)
    }

    // Update the datagram's signature field with the generated signature
    copy(dg.Signature[:], []byte(signature)) // Ensure we copy the signature into the byte array

    // Return the serialized data including the signature
    return serializedData, nil
}
