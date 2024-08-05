package handlers

import (
    "fmt"
    "path/filepath"
    "resilience/main"
)

// sendErrorResponse prepares and sends an error response datagram.
func sendErrorResponse(ctx main.HandlerContext, errorMessage string) error {
    var responseDg main.ResponseDatagram
    // Use the original signature as the nonce, dereferencing ctx.Datagram
    copy(responseDg.Nonce[:], ctx.Datagram.Signature[:]) 
    responseDg.Result[0] = 1 // Set error code
    copy(responseDg.Result[1:], []byte(errorMessage)) // Copy the error message

    // Construct the account directory path directly, dereferencing ctx.Datagram
    accountDir := filepath.Join(main.Datadir, "accounts", string(ctx.Datagram.XUsername[:]))
    
    // Generate signature for response datagram
    if err := main.SignResponseDatagram(&responseDg, accountDir); err != nil {
        fmt.Printf("Failed to sign response datagram: %v\n", err)
        return err
    }

    // Send the signed response datagram
    _, err := ctx.Conn.WriteToUDP(responseDg[:], ctx.Addr)
    if err != nil {
        return fmt.Errorf("error sending response datagram: %w", err)
    }

    return nil
}
