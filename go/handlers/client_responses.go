package handlers

import (
    "fmt"
    "net"
    "resilience/main"
)

// sendErrorResponse prepares and sends an error response datagram.
func sendErrorResponse(dg main.Datagram, addr *net.UDPAddr, conn *net.UDPConn, accountDir string, errorMessage string) error {
    var responseDg main.ResponseDatagram
    copy(responseDg.Nonce[:], dg.Signature[:]) // Use the original signature as the nonce
    responseDg.Result[0] = 1                    // Set error code
    copy(responseDg.Result[1:], []byte(errorMessage)) // Copy the error message

    // Generate signature for response datagram
    if err := main.SignResponseDatagram(&responseDg, accountDir); err != nil {
        fmt.Printf("Failed to sign response datagram: %v\n", err)
        return err
    }

    // Send the signed response datagram
    _, err := conn.WriteToUDP(responseDg[:], addr)
    if err != nil {
        return fmt.Errorf("error sending response datagram: %w", err)
    }

    return nil
}
