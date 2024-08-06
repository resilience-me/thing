package handlers

import (
    "fmt"
    "resilience/main"
)

// prepareAndSendResponse handles the common tasks for sending a response datagram.
func prepareAndSendResponse(ctx main.HandlerContext, resultCode byte, message []byte) error {
    var responseDg main.ResponseDatagram
    copy(responseDg.Nonce[:], ctx.Datagram.Signature[:])
    responseDg.Result[0] = resultCode // Set result code
    copy(responseDg.Result[1:], message) // Copy the message or data

    // Generate signature for response datagram
    username := string(ctx.Datagram.XUsername[:])
    if err := main.SignResponseDatagram(&responseDg, username); err != nil {
        return fmt.Errorf("failed to sign response datagram: %w", err)
    }

    // Send the signed response datagram
    _, err := ctx.Conn.WriteToUDP(responseDg[:], ctx.Addr)
    if err != nil {
        return fmt.Errorf("error sending response datagram: %w", err)
    }

    return nil
}

// SendErrorResponse prepares and sends an error response datagram.
func SendErrorResponse(ctx main.HandlerContext, errorMessage string) error {
    return prepareAndSendResponse(ctx, 1, []byte(errorMessage)) // Error code 1
}

// SendSuccessResponse prepares and sends a success response datagram.
func SendSuccessResponse(ctx main.HandlerContext, resultData []byte) error {
    return prepareAndSendResponse(ctx, 0, resultData) // Success code 0
}
