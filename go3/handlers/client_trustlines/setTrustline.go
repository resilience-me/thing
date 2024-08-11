package client_trustlines

import (
    "encoding/binary"
    "fmt"
    "ripple/database/db_trustlines"
    "ripple/handlers" // Import the handlers package
)

// Assuming the Session and Datagram structures are defined as per your latest setup
func SetTrustline(session Session) {
    // We assume session.Datagram is directly accessible and correctly initialized
    datagram := session.Datagram

    // Retrieve the previous counter value using the getter
    prevCounter, err := db_trustlines.GetCounter(datagram)
    if err != nil {
        fmt.Printf("Error getting previous counter: %v\n", err) // Log detailed error
        _ = handlers.SendErrorResponse("Failed to read counter file.", session.Conn) // Send simpler error message
        return
    }

    // Check the counter directly as it is already a uint32
    if datagram.Counter <= prevCounter {
        fmt.Println("Received counter is not greater than previous counter. Potential replay attack.")
        _ = handlers.SendErrorResponse("Received counter is not valid.", session.Conn) // Send simpler error message
        return
    }

    // Retrieve the trustline amount from the Datagram
    trustlineAmount := binary.BigEndian.Uint32(datagram.Arguments[:4])

    // Write the new trustline amount using the setter
    if err := db_trustlines.SetTrustlineOut(datagram, trustlineAmount); err != nil {
        fmt.Printf("Error writing trustline to file: %v\n", err) // Log detailed error
        _ = handlers.SendErrorResponse("Failed to write trustline.", session.Conn) // Send simpler error message
        return
    }

    // Write the new counter value using the setter
    if err := db_trustlines.SetCounter(datagram, datagram.Counter); err != nil {
        fmt.Printf("Error writing counter to file: %v\n", err) // Log detailed error
        _ = handlers.SendErrorResponse("Failed to write counter.", session.Conn) // Send simpler error message
        return
    }

    fmt.Println("Trustline and counter updated successfully.")

    if err := handlers.SendSuccessResponse("Trustline updated successfully."); err != nil {
        fmt.Printf("Failed to send success response: %v\n", err) // Log detailed error if any
        return
    }

    fmt.Println("Sent success response to client.")
}
