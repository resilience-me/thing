package server_trustlines

import (
    "fmt"
    "time"
    "ripple/main"                    // Updated to match your import structure
    "ripple/handlers"                // Updated to match your import structure
    "ripple/database/db_trustlines"  // Updated to match your import structure
)

// SetTrustline handles setting or updating a trustline from another server's perspective
func SetTrustline(session main.Session) {
    // Retrieve the counter_in value using the new getter
    counterIn, err := db_trustlines.GetCounterIn(session.Datagram)
    if err != nil {
        fmt.Printf("Error getting counter_in: %v\n", err)
        return
    }

    // Check the counter
    counter := session.Datagram.Counter // Directly using the uint32 counter
    if counter <= counterIn {
        fmt.Println("Received counter is not greater than counter_in. Potential replay attack.")
        return
    }

    // Retrieve the trustline amount from the Datagram
    trustlineAmount := main.BytesToUint32(session.Datagram.Arguments[:4])

    // Retrieve the sync_in value using the new getter
    prevSyncIn, err := db_trustlines.GetSyncIn(session.Datagram)
    if err != nil {
        fmt.Printf("Error getting counter_in: %v\n", err)
        return
    }

    // Retrieve the syncIn counter from the Datagram
    syncInBytes := session.Datagram.Arguments[4:8]
    syncIn := main.BytesToUint32(syncInBytes)
    if syncIn <= prevSyncIn {
        fmt.Println("The local sync_in value is synchronized with the peers most recent trustline_out update.")
    } else {
        // Write the new trustline amount using the setter
        if err := db_trustlines.SetTrustlineIn(session.Datagram, trustlineAmount); err != nil {
            fmt.Printf("Error writing trustline to file: %v\n", err)
            return
        }
    
        // Write the new sync_in value using the setter
        if err := db_trustlines.SetSyncIn(session.Datagram, syncIn); err != nil {
            fmt.Printf("Error writing sync_in to file: %v\n", err)
            return
        }
    
        fmt.Println("trustline_in and sync_in updated successfully.")
    }

    // Write the Unix timestamp using the setter
    if err := db_trustlines.SetTimestamp(session.Datagram, time.Now().Unix()); err != nil {
        fmt.Printf("Error writing timestamp to file: %v\n", err)
        return
    }

    if syncIn > prevSyncIn {

        // Retrieve the counter_in value using the new getter
        counterOut, err := db_trustlines.GetCounterOut(session.Datagram)
        if err != nil {
            fmt.Printf("Error getting counter_out: %v\n", err)
            return
        }

        // Prepare the datagram to send back to the peer
        dg := main.Datagram{
            Command:        main.ServerTrustlines_SetSyncOut,
            Username:      session.Datagram.PeerUsername,       // Reverse the usernames for response
            PeerUsername:      session.Datagram.Username,
            PeerServerAddress: main.GetServerAddress(),      // Use the server's address
            Arguments:     syncInBytes,
            Counter:        counterOut,
        }
    
        buf := handlers.CreateSignedDatagram(session, &dg)

        // Replace explicit signing and sending with the centralized function call
        if err := handlers.SendDatagram(session, &dg); err != nil {
            fmt.Printf("Failed to sign and send datagram: %v\n", err)
            return
        }
    
        // Add a success message indicating all operations were successful
        fmt.Println("Trustline update and datagram sending completed successfully.")
    }
}
