package trustlines

import (
    "encoding/binary"
    "fmt"
    "ripple/main"
    "ripple/database/db_trustlines"
)

// PrepareDatagramForSync prepares the datagram for synchronization based on the sync status.
func PrepareDatagramForSync(datagram *main.Datagram) (*main.Datagram, uint32, error) {
    // Retrieve the syncCounter and sync status
    syncCounter, isSynced, err := GetSyncStatus(datagram)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to retrieve sync status for user %s: %v", datagram.Username, err)
    }

    // Retrieve and increment the counter_out value
    counterOut, err := GetAndIncrementCounterOut(datagram)
    if err != nil {
        return nil, 0, fmt.Errorf("error handling counter_out for user %s: %v", datagram.Username, err)
    }

    dg := &main.Datagram{
        Username:          datagram.PeerUsername,
        PeerUsername:      datagram.Username,
        PeerServerAddress: datagram.PeerServerAddress,
        Counter:           counterOut,
    }

    if isSynced {
        // Trustline is already synced, so prepare a SetTimestamp command
        dg.Command = main.ServerTrustlines_SetTimestamp
    } else {
        // Trustline is not synced, prepare to send the trustline
        trustline, err := db_trustlines.GetTrustlineOut(datagram)
        if err != nil {
            return nil, 0, fmt.Errorf("error getting trustline for user %s: %v", datagram.Username, err)
        }
        dg.Command = main.ServerTrustlines_SetTrustline
        binary.BigEndian.PutUint32(dg.Arguments[:4], trustline)
        binary.BigEndian.PutUint32(dg.Arguments[4:8], syncCounter)
    }

    return dg, syncCounter, nil
}
