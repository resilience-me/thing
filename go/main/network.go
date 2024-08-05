// SignAndSendResponse signs and sends the response datagram.
func SignAndSendResponseDatagram(responseDg *ResponseDatagram, addr *net.UDPAddr, conn *net.UDPConn, dir string) error {
    // Call generateSignature directly with the ResponseDatagram's byte representation
    signature, err := generateSignature((*responseDg)[:], dir)
    if err != nil {
        return fmt.Errorf("failed to generate signature for ResponseDatagram: %w", err)
    }

    // Copy the generated signature into the response datagram's signature field
    copy(responseDg.Signature[:], signature)

    // Send the signed response datagram
    _, err := conn.WriteToUDP(responseDg[:], addr)
    if err != nil {
        return fmt.Errorf("error sending response datagram: %w", err)
    }

    return nil
}

// SignAndSendResponse signs and sends the response datagram.
func SignAndSendDatagram(dg *ResponseDatagram, addr *net.UDPAddr, conn *net.UDPConn, dir string) error {
    // Call generateSignature directly with the Datagram's byte representation
    signature, err := generateSignature((*dg)[:], dir)
    if err != nil {
        return fmt.Errorf("failed to generate signature for Datagram: %w", err)
    }

    // Copy the generated signature into the datagram's signature field
    copy(dg.Signature[:], signature)

    // Send the signed response datagram
    _, err := conn.WriteToUDP(responseDg[:], addr)
    if err != nil {
        return fmt.Errorf("error sending response datagram: %w", err)
    }

    return nil
}
