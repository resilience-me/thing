package handlers

// InitializeDatagram initializes common Datagram fields for a response and increments counter_out.
func InitializeDatagram(datagram *types.Datagram) (*types.Datagram, error) {
    // Retrieve and increment the counter_out value
    counterOut, err := db_server.GetAndIncrementCounterOut(datagram)
    if err != nil {
        return nil, fmt.Errorf("error handling counter_out for user %s: %v", datagram.Username, err)
    }

    dg := types.NewDatagram(datagram.PeerUsername, datagram.Username, counterOut)

    return dg, nil
}
