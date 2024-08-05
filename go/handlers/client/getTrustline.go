// package client

// import (
//     "encoding/binary"
//     "fmt"
//     "net"
//     "os"
//     "path/filepath"
//     "strconv"

//     "resilience/main"
//     "resilience/handlers"
// )

// // getTrustline handles fetching the trustline information for both inbound and outbound.
// func getTrustline(ctx main.HandlerContext, filename string) {
//     // Validate the client request
//     if err := handlers.ValidateClientRequest(ctx); err != nil {
//         fmt.Printf("Validation failed: %v\n", err) // Log detailed error
//         return // Error response has already been sent in ValidateClientRequest
//     }

//     // Get the peer directory for the trustline
//     peerDir := main.GetPeerDir(ctx.Datagram)

//     trustlinePath := filepath.Join(peerDir, "trustline", filename)
//     trustlineAmountStr, err := os.ReadFile(trustlinePath)
//     if err != nil {
//         fmt.Printf("Error reading trustline file (%s): %v\n", filename, err) // Log the error
//         _ = handlers.SendErrorResponse(ctx, "Error reading trustline file.")
//         return
//     }

//     // Convert the string to an integer
//     trustlineAmount, err := strconv.ParseUint(string(trustlineAmountStr), 10, 32)
//     if err != nil {
//         fmt.Printf("Error converting trustline amount to integer: %v\n", err) // Log the error
//         _ = handlers.SendErrorResponse(ctx, "Error converting trustline amount to integer.")
//         return
//     }

//     // Prepare success response
//     responseData := make([]byte, 4) // Allocate 4 bytes for the trustline amount
//     binary.BigEndian.PutUint32(responseData, uint32(trustlineAmount)) // Convert the trustline amount to bytes

//     // Send the success response back to the client
//     if err := handlers.SendSuccessResponse(ctx, responseData); err != nil {
//         fmt.Printf("Error sending success response: %v\n", err) // Log the error
//         return
//     }

//     fmt.Printf("Trustline amount (%s) sent successfully.\n", filename)
// }

// // GetTrustlineIn handles fetching the inbound trustline information
// func GetTrustlineIn(ctx main.HandlerContext) {
//     getTrustline(ctx, "trustline_in.txt")
// }

// // GetTrustlineOut handles fetching the outbound trustline information
// func GetTrustlineOut(ctx main.HandlerContext) {
//     getTrustline(ctx, "trustline_out.txt")
// }
