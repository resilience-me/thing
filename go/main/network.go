package main

import (
    "fmt"
    "net"
)

// SignAndSendResponseDatagram signs and sends the response datagram.
func SignAndSendResponseDatagram(responseDg *ResponseDatagram, addr *net.UDPAddr, conn *net.UDPConn, accountDir string) error {
    // Generate signature for response datagram
    if err := main.SignResponseDatagram(&responseDg, accountDir); err != nil {
        fmt.Printf("Failed to sign response datagram: %v\n", err)
        return
    }
    // Send the signed response datagram
    _, err = conn.WriteToUDP(responseDg[:], addr)
    if err != nil {
        return fmt.Errorf("error sending response datagram: %w", err)
    }

    return nil
}

// SignAndSendDatagram signs and sends the datagram
func SignAndSendDatagram(dg *Datagram, addr *net.UDPAddr, conn *net.UDPConn, peerDir string) error {
    // Generate signature for Datagram
    if err := main.SignDatagram(&dg, peerDir); err != nil {
        fmt.Printf("Failed to sign response datagram: %v\n", err)
        return
    }

    // Send the signed datagram
    _, err := conn.WriteToUDP(dg[:], addr)
    if err != nil {
        return fmt.Errorf("error sending datagram: %w", err)
    }

    return nil
}
