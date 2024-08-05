package main

import (
    "fmt"
    "net"
)

// SignAndSendResponseDatagram signs and sends the response datagram.
func SignAndSendResponseDatagram(responseDg *ResponseDatagram, addr *net.UDPAddr, conn *net.UDPConn, dir string) error {
    // Generate signature for ResponseDatagram
    signature, err := GenerateSignature((*responseDg)[:], dir)
    if err != nil {
        return fmt.Errorf("failed to generate signature for ResponseDatagram: %w", err)
    }

    // Copy the generated signature into the response datagram's signature field
    copy(responseDg.Signature[:], signature)

    // Send the signed response datagram
    _, err = conn.WriteToUDP(responseDg[:], addr)
    if err != nil {
        return fmt.Errorf("error sending response datagram: %w", err)
    }

    return nil
}

// SignAndSendDatagram signs and sends the datagram
func SignAndSendDatagram(dg *Datagram, addr *net.UDPAddr, conn *net.UDPConn, dir string) error {
    // Generate signature for Datagram
    signature, err := GenerateSignature((*dg)[:], dir)
    if err != nil {
        return fmt.Errorf("failed to generate signature for Datagram: %w", err)
    }

    // Copy the generated signature into the datagram's signature field
    copy(dg.Signature[:], signature)

    // Send the signed datagram
    _, err := conn.WriteToUDP(dg[:], addr)
    if err != nil {
        return fmt.Errorf("error sending datagram: %w", err)
    }

    return nil
}
