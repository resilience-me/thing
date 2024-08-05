package handlers

import (
    "fmt"
    "resilience/main"
)

// sendErrorResponse prepares and sends an error response datagram.
func SendErrorResponse(ctx main.HandlerContext, errorMessage string) error {
    var responseDg main.ResponseDatagram
    // Use the original signature as the nonce, dereferencing ctx.Datagram
    copy(responseDg.Nonce[:], ctx.Datagram.Signature[:]) 
    responseDg.Result[0] = 1 // Set error code
    copy(responseDg.Result[1:], []byte(errorMessage)) // Copy the error message

    // Use GetAccountDir to construct the account directory path
    accountDir := main.GetAccountDir(ctx.Datagram)
    
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

func SendSuccessResponse(ctx main.HandlerContext, result []byte) error {
    var responseDg main.ResponseDatagram
    copy(responseDg.Nonce[:], ctx.Datagram.Signature[:])
    responseDg.Result[0] = 0 // Set success code
    copy(responseDg.Result[1:], result)

    accountDir := main.GetAccountDir(ctx.Datagram)
    if err := main.SignResponseDatagram(&responseDg, accountDir); err != nil {
        fmt.Printf("Failed to sign response datagram: %v\n", err)
        return err
    }

    _, err := ctx.Conn.WriteToUDP(responseDg[:], ctx.Addr)
    if err != nil {
        return fmt.Errorf("error sending response datagram: %w", err)
    }

    return nil
}
