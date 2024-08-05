package client

import (
    "fmt"
    "net"
    "os"
    "path/filepath"
    "strconv"

    "resilience/main"
    "resilience/handlers"
)

// Handles fetching the inbound trustline information
func GetTrustlineIn(ctx main.HandlerContext) {
    // First, get the account directory
    accountDir, err := main.GetAccountDir(ctx.Datagram) // No need for the & operator
    if err != nil {
        fmt.Printf("Error getting account directory: %v\n", err) // Log the error
        _ = handlers.SendErrorResponse(ctx, "Error getting account directory.")
        return
    }

    // Now, get the peer directory using the account directory
    peerDir, err := main.GetPeerDir(ctx.Datagram, accountDir) // Pass accountDir
    if err != nil {
        fmt.Printf("Error getting peer directory: %v\n", err) // Log the error
        _ = handlers.SendErrorResponse(ctx, "Error getting peer directory for inbound trustline.")
        return
    }

    trustlineInPath := filepath.Join(peerDir, "trustline", "trustline_in.txt")
    trustlineAmountStr, err := os.ReadFile(trustlineInPath)
    if err != nil {
        fmt.Printf("Error reading inbound trustline file: %v\n", err) // Log the error
        _ = handlers.SendErrorResponse(ctx, "Error reading inbound trustline file.")
        return
    }

    // Convert the string to an integer
    trustlineAmount, err := strconv.ParseUint(string(trustlineAmountStr), 10, 32)
    if err != nil {
        fmt.Printf("Error converting trustline amount to integer: %v\n", err) // Log the error
        _ = handlers.SendErrorResponse(ctx, "Error converting trustline amount to integer.")
        return
    }

    // Prepare response datagram
    var responseDg main.ResponseDatagram
    responseDg.Result[0] = 0 // Set success code
    copy(responseDg.Nonce[:], ctx.Datagram.Signature[:]) // Use the original signature as the nonce

    // Store the trustline amount as bytes in the response
    binary.BigEndian.PutUint32(responseDg.Result[1:], uint32(trustlineAmount)) // Convert back to bytes
    if err := main.SignResponseDatagram(&responseDg, accountDir); err != nil {
        fmt.Printf("Failed to sign response datagram: %v\n", err) // Log the error
        _ = handlers.SendErrorResponse(ctx, "Failed to sign response datagram.")
        return
    }

    // Send the response back to the client
    _, err = ctx.Conn.WriteToUDP(responseDg[:], ctx.Addr)
    if err != nil {
        fmt.Printf("Error sending inbound trustline amount: %v\n", err) // Log the error
        return
    }

    fmt.Println("Inbound trustline amount sent successfully.")
}
