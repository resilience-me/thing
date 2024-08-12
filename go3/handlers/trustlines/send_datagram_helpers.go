package trustlines

import (
    "ripple/main"
    "ripple/database/db_trustlines"
    "fmt"
)

// InitializeDatagram initializes common Datagram fields for a response and increments counter_out.
func InitializeDatagram(datagram *main.Datagram) (*main.Datagram, error) {
    // Retrieve and increment the counter_out value
    counterOut, err := GetAndIncrementCounterOut(datagram)
    if err != nil {
        return nil, fmt.Errorf("error handling counter_out for user %s: %v", datagram.Username, err)
    }

    dg := &main.Datagram{
        Username:          datagram.PeerUsername,
        PeerUsername:      datagram.Username,
        PeerServerAddress: main.GetServerAddress(),
        Counter:           counterOut,
    }

    return dg, nil
}
