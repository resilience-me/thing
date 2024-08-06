package handlers

import (
    "fmt"
    "resilience/main"
)

// prepareAndSendResponse is a shared function to set up and send a response datagram.
func prepareAndSendResponse(ctx main.HandlerContext, resultCode byte, message []byte) error {
    var responseDg main.ResponseDatagram
    // Use the original signature as the nonce, dereferencing ctx.Datagram
    copy(responseDg.Nonce[:], ctx.Datagram.Signature[:])
    responseDg.Result[0] = resultCode // Set result code (error or success)
    copy(responseDg.Result[1:], message) // Copy the message or data

    // Generate signature for response datagram
    username := string(ctx.Datagram.XUsername[:])
    if err := main.SignResponseDatagram(&responseDg, username); err != nil {
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

// SendErrorResponse prepares and sends an error response datagram using the shared function.
func SendErrorResponse(ctx main.HandlerContext, errorMessage string) error {
    return prepareAndSendResponse(ctx, 1, []byte(errorMessage)) // Set error code to 1
}

// SendSuccessResponse prepares and sends a success response datagram using the shared function.
func SendSuccessResponse(ctx main.HandlerContext, resultData []byte) error {
    return prepareAndSendResponse(ctx, 0, resultData) // Set success code to 0
}
