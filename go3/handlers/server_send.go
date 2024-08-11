package handlers

import (
    "fmt"
    "net"
    "ripple/main" // Assuming this contains your Datagram and Session types
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


// SendDatagram sends a datagram to the specified server address and port.
func SendDatagram(session main.Session, data []byte) error {
    // Get the server address from the session
    serverAddress := session.Datagram.PeerServerAddress

    // Resolve the address, treating it as an IP address or domain name
    addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", serverAddress, 2012))
    if err != nil {
        return fmt.Errorf("failed to resolve server address '%s': %w", serverAddress, err)
    }

    // Establish a connection to the server
    conn, err := net.DialTCP("tcp", nil, addr)
    if err != nil {
        return fmt.Errorf("failed to connect to server '%s': %w", serverAddress, err)
    }
    defer conn.Close() // Ensure the connection is closed after sending

    // Send the data to the server
    if _, err := conn.Write(data); err != nil {
        return fmt.Errorf("failed to send data to server '%s': %w", serverAddress, err)
    }

    return nil // Successfully sent the datagram
}

// SignAndSendDatagram creates a signed datagram and sends it over the network.
func SignAndSendDatagram(session main.Session, dg *main.Datagram) error {
    // Create the signed datagram
    serializedData, err := CreateSignedDatagram(session, dg)
    if err != nil {
        return fmt.Errorf("failed to create signed datagram: %w", err)
    }

    // Send the signed datagram over the network
    if err := SendDatagram(session, serializedData); err != nil {
        return fmt.Errorf("failed to send datagram: %w", err)
    }

    return nil // Successfully signed and sent
}
