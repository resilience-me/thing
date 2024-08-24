package handlers

import (
    "fmt"
    "ripple/auth"
    "ripple/comm"
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
    datagram, err := PrepareDatagramWithoutCommand(username, peerServerAddress, peerUsername)
    if err != nil {
        return nil, fmt.Errorf("Failed to prepare datagram: %v", err)
    }
    datagram.Command = command
    copy(datagram.Arguments[:], arguments)

    return datagram, nil
}

// PrepareDatagramResponse calls PrepareDatagram with fields from an incoming datagram
func PrepareDatagramResponse(datagram *types.Datagram) (*types.Datagram, error) {
    return PrepareDatagramWithoutCommand(datagram.Username, datagram.PeerServerAddress, datagram.PeerUsername)
}

// PrepareAndSendDatagram prepares, signs, and sends a datagram to a specified peer.
func PrepareAndSendDatagram(command byte, username, serverAddress, peerUsername string, arguments []byte) error {
    // Prepare the datagram with the command and arguments
    newDatagram, err := PrepareDatagram(command, username, serverAddress, peerUsername, arguments)
    if err != nil {
        return fmt.Errorf("Failed to prepare datagram: %v", err)
    }

    // Sign and send the datagram to the target peer
    if err := comm.SignAndSendDatagram(newDatagram, serverAddress); err != nil {
        return fmt.Errorf("Failed to sign and send datagram to %s at %s: %v", peerUsername, serverAddress, err)
    }

    return nil
}
