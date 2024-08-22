package comm

import (
    "fmt"
    "ripple/types"
    "ripple/auth"
)

// signAndSendDatagram creates a signed datagram and sends it over the network with custom priority
func signAndSendDatagram(dg *types.Datagram, peerServerAddress string, maxRetries int) error {
    // Create the signed datagram
    serializedData, err := auth.SignDatagram(dg, peerServerAddress)
    if err != nil {
        return fmt.Errorf("failed to create signed datagram: %w", err)
    }
    
    // Send the signed datagram over the network
    if err := SendWithResolvedAddress(peerServerAddress, serializedData, maxRetries); err != nil {
        return fmt.Errorf("failed to send datagram: %w", err)
    }

    return nil // Successfully signed and sent
}

// SignAndSendDatagram creates a signed datagram and sends it over the network with low priority.
func SignAndSendDatagram(dg *types.Datagram, peerServerAddress string) error {
    return signAndSendDatagram(dg, peerServerAddress, LowImportance)
}

// SignAndSendPriorityDatagram creates a signed datagram and sends it over the network with high priority.
func SignAndSendPriorityDatagram(dg *types.Datagram, peerServerAddress string) error {
    return signAndSendDatagram(dg, peerServerAddress, HighImportance)
}
