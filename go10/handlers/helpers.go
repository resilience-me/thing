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
