package handlers

import (
    "ripple/types"
    "ripple/auth"
    "ripple/types"

)

// PrepareDatagram prepares common Datagram fields and increments counter_out.
func PrepareDatagram(username, peerServerAddress, peerUsername string) (*types.Datagram, error) {
    // Retrieve and increment the counter_out value
    counterOut, err := auth.GetAndIncrementCounterOut(username, peerServerAddress, peerUsername)
    if err != nil {
        return nil, fmt.Errorf("error handling counter_out for user %s: %v", username, err)
    }

    dg := types.NewDatagram(peerUsername, username, counterOut)

    return dg, nil
}

// PrepareDatagramResponse calls PrepareDatagram with fields from an incoming datagram
func PrepareDatagramResponse(dg *types.Datagram) (*types.Datagram, error) {
    return PrepareDatagram(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
}

// signAndSendDatagram creates a signed datagram and sends it over the network.
func signAndSendDatagram(dg *types.Datagram, peerServerAddress string, maxRetries int) error {
    // Create the signed datagram
    serializedData, err := auth.SignDatagram(dg, peerServerAddress)
    if err != nil {
        return fmt.Errorf("failed to create signed datagram: %w", err)
    }
    
    // Send the signed datagram over the network
    if err := comm.SendWithResolvedAddress(peerServerAddress, serializedData, maxRetries); err != nil {
        return fmt.Errorf("failed to send datagram: %w", err)
    }

    return nil // Successfully signed and sent
}

// SignAndSendDatagram creates a signed datagram and sends it over the network.
func SignAndSendDatagram(dg *types.Datagram, peerServerAddress string) error {
    return signAndSendDatagram(dg, peerServerAddress, comm.LowImportance)
}

// SignAndSendDatagram creates a signed datagram and sends it over the network.
func SignAndSendPriorityDatagram(dg *types.Datagram, peerServerAddress string) error {
    return signAndSendDatagram(dg, peerServerAddress, comm.HighImportance)
}
