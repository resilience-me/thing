// SignAndSendDatagram signs and sends a Datagram over UDP.
func SignAndSendDatagram(ctx main.HandlerContext, dg *main.Datagram) error {
    // Sign the datagram
    if err := main.SignDatagram(dg); err != nil {
        return fmt.Errorf("Failed to sign datagram: %v\n", err)
    }

    // Send the datagram back to the peer
    _, err := ctx.Conn.WriteToUDP(dg[:], ctx.Addr)
    if err != nil {
        fmt.Printf("Error sending datagram: %v\n", err)
        return err
    }

    return nil
}
