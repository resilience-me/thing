package auth

import (
    "fmt"
    "ripple/types"
)

// SignDatagram creates a signed datagram by serializing it and adding a signature.
// It requires the session to load the secret key for signature generation.
func SignDatagram(dg *types.Datagram, peerServerAddress string) ([]byte, error) {
    // Serialize the datagram without the signature field
    serializedData, err := types.SerializeDatagram(dg)
    if err != nil {
        return nil, fmt.Errorf("failed to serialize datagram: %w", err)
    }

    // Load the secret key for signature generation
    secretKey, err := loadServerSecretKeyOut(dg, peerServerAddress)
    if err != nil {
        return nil, fmt.Errorf("failed to load server secret key: %w", err)
    }

    // Generate signature for the serialized data
    signature := generateSignature(serializedData[:357], secretKey)

    // Update the datagram's signature field with the generated signature
    copy(dg.Signature[:], []byte(signature)) // Ensure we copy the signature into the byte array

    // Return the serialized data including the signature
    return serializedData, nil
}
