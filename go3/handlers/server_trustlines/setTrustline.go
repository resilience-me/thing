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
    // Retrieve the sync_in value using the new getter
    syncIn, err := db_trustlines.GetSyncIn(session.Datagram)
    if err != nil {
        fmt.Printf("Error getting sync_in: %v\n", err)
        return
    }

    // Check the counter
    counter := session.Datagram.Counter // Directly using the uint32 counter
    if counter <= syncIn {
        fmt.Println("Received counter is not greater than sync_in. Potential replay attack.")
        return
    }

    // Retrieve the trustline amount from the Datagram
    trustlineAmount := binary.BigEndian.Uint32(session.Datagram.Arguments[:4])

    // Write the new trustline amount using the setter
    if err := db_trustlines.SetTrustlineIn(session.Datagram, trustlineAmount); err != nil {
        fmt.Printf("Error writing trustline to file: %v\n", err)
        return
    }

    // Write the new sync_in value using the setter
    if err := db_trustlines.SetSyncIn(session.Datagram, counter); err != nil {
        fmt.Printf("Error writing sync_in to file: %v\n", err)
        return
    }

    // Write the Unix timestamp using the setter
    if err := db_trustlines.SetTimestamp(session.Datagram, time.Now().Unix()); err != nil {
        fmt.Printf("Error writing timestamp to file: %v\n", err)
        return
    }

    fmt.Println("Trustline, sync_in, and timestamp updated successfully.")

    // Prepare the datagram to send back to the peer
    dg := main.Datagram{
        Command:        main.ServerTrustlines_SetSyncOut,
        XUsername:      session.Datagram.YUsername,       // Reverse the usernames for response
        YUsername:      session.Datagram.XUsername,
        YServerAddress: main.GetServerAddress(),      // Use the server's address
        Counter:        counter,                        // Directly using the uint32 counter
    }

    // Replace explicit signing and sending with the centralized function call
    if err := handlers.SignAndSendDatagram(session, &dg); err != nil {
        fmt.Printf("Failed to sign and send datagram: %v\n", err)
        return
    }

    // Add a success message indicating all operations were successful
    fmt.Println("Trustline update and datagram sending completed successfully.")
}
