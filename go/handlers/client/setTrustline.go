// SetTrustline handles setting or updating a trustline from the client's perspective
func SetTrustline(ctx main.HandlerContext) {
    trustlineAmount := binary.BigEndian.Uint32(ctx.Datagram.Arguments[:4])

    accountDir, err := main.GetAccountDir(ctx.Datagram)
    if err != nil {
        fmt.Printf("Error getting account directory: %v\n", err) // Log detailed error
        sendErrorResponse(ctx, "Failed to get account directory.") // Send simpler error message
        return
    }

    peerDir, err := main.GetPeerDir(ctx.Datagram, accountDir)
    if err != nil {
        fmt.Printf("Error getting peer directory: %v\n", err) // Log detailed error
        sendErrorResponse(ctx, "Failed to get peer directory.") // Send simpler error message
        return
    }

    if err := main.VerifySignature(ctx.Datagram, accountDir); err != nil {
        fmt.Printf("Signature verification failed: %v\n", err) // Log detailed error
        sendErrorResponse(ctx, "Signature verification failed.") // Send simpler error message
        return
    }

    // Get the trustline directory
    trustlineDir := filepath.Join(peerDir, "trustline")

    // Construct the trustline and counter file paths
    counterOutPath := filepath.Join(trustlineDir, "counter_out.txt")
    trustlineOutPath := filepath.Join(trustlineDir, "trustline_out.txt")

    // Load the previous counter value
    prevCounterStr, err := os.ReadFile(counterOutPath)
    if err != nil && !os.IsNotExist(err) {
        fmt.Printf("Error reading counter file: %v\n", err) // Log detailed error
        sendErrorResponse(ctx, "Failed to read counter file.") // Send simpler error message
        return
    }

    prevCounter, err := strconv.ParseUint(string(prevCounterStr), 10, 32) // Parse as uint64 first
    if err != nil {
        fmt.Printf("Error parsing previous counter string: %v\n", err) // Log detailed error
        sendErrorResponse(ctx, "Failed to parse previous counter.") // Send simpler error message
        return
    }

    // Check the counter
    counter := binary.BigEndian.Uint32(ctx.Datagram.Counter[:])
    if counter <= uint32(prevCounter) {
        fmt.Println("Received counter is not greater than previous counter. Potential replay attack.")
        sendErrorResponse(ctx, "Received counter is not valid.") // Send simpler error message
        return
    }

    // Write the new trustline amount to the file
    if err := os.WriteFile(trustlineOutPath, []byte(fmt.Sprintf("%d", trustlineAmount)), 0644); err != nil {
        fmt.Printf("Error writing trustline to file: %v\n", err) // Log detailed error
        sendErrorResponse(ctx, "Failed to write trustline.") // Send simpler error message
        return
    }

    // Write the new counter value as a string
    counterStr := fmt.Sprintf("%d", counter)
    if err := os.WriteFile(counterOutPath, []byte(counterStr), 0644); err != nil {
        fmt.Printf("Error writing counter to file: %v\n", err) // Log detailed error
        sendErrorResponse(ctx, "Failed to write counter.") // Send simpler error message
        return
    }

    fmt.Println("Trustline and counter updated successfully.")

    // Prepare response datagram
    var responseDg main.ResponseDatagram
    responseDg.Result[0] = 0 // Set success code
    copy(responseDg.Nonce[:], ctx.Datagram.Signature[:]) // Use the original signature as the nonce
    copy(responseDg.Result[1:], []byte("Trustline updated successfully.")) // More informative success message

    if err := main.SignResponseDatagram(&responseDg, accountDir); err != nil {
        fmt.Printf("Failed to sign response datagram: %v\n", err) // Log detailed error
        sendErrorResponse(ctx, "Failed to sign response datagram.") // Send simpler error message
        return
    }

    // Send the response back to the client
    _, err = ctx.Conn.WriteToUDP(responseDg[:], ctx.Addr)
    if err != nil {
        fmt.Printf("Error sending response to client: %v\n", err) // Log detailed error
        return
    }
    fmt.Println("Sent success response to client.")
}
