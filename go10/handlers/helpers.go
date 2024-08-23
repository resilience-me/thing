package handlers

import (
    "ripple/types"
    "ripple/auth"
    "ripple/types"

)

// PrepareDatagramWithoutCommand prepares common Datagram fields and increments counter_out.
func PrepareDatagramWithoutCommand(username, peerServerAddress, peerUsername string) (*types.Datagram, error) {
    // Retrieve and increment the counter_out value
    counterOut, err := auth.GetAndIncrementCounterOut(username, peerServerAddress, peerUsername)
    if err != nil {
        return nil, fmt.Errorf("error handling counter_out for user %s: %v", username, err)
    }

    dg := types.NewDatagram(peerUsername, username, counterOut)

    return dg, nil
}

// PrepareDatagram prepares a datagram with all necessary fields including the command and arguments.
func PrepareDatagram(command byte, username, peerServerAddress, peerUsername string, arguments []byte) (*types.Datagram, error) {
    // Prepare the new datagram
    datagram, err := handlers.PrepareDatagramWithoutCommand(datagram.Username, peer.ServerAddress, peer.Username)
    if err != nil {
        return nil, fmt.Errorf("Failed to prepare datagram: %v", err)
    }
    datagram.Command = command
    copy(dg.Arguments[:], arguments)

    return datagram, nil
}

// PrepareDatagramResponse calls PrepareDatagram with fields from an incoming datagram
func PrepareDatagramResponse(dg *types.Datagram) (*types.Datagram, error) {
    return PrepareDatagram(dg.Username, dg.PeerServerAddress, dg.PeerUsername)
}
