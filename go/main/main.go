package main

import (
    "fmt"
    "net"
)

func main() {
    if err := initConfig(); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to initialize configuration: %v\n", err)
        return
    }

    addr := net.UDPAddr{
        Port: 2012,
        IP:   net.ParseIP("::"), // Listen on all IPv6 and mapped IPv4 addresses.
    }

    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error listening on UDP port %d: %v\n", addr.Port, err)
        return
    }
    defer conn.Close()

    fmt.Printf("Server is listening on all network interfaces for both IPv4 and IPv6 at port %d\n", addr.Port)

    for {
        var dg Datagram
        n, remoteAddr, err := conn.ReadFromUDP(dg[:])
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error reading from UDP: %v\n", err)
            continue
        }
        if n != len(dg) {
            fmt.Printf("Received incorrect datagram size from %s. Expected %d bytes, got %d bytes.\n", remoteAddr, len(dg), n)
            continue
        }

        // Determine if the received command is a client command (where the most significant bit is 0)
        isClientCommand := dg.Command < 128

        // Check if the source IP address is not a loopback address (i.e., not localhost)
        isNotLoopback := !remoteAddr.IP.IsLoopback()

        // If LocalClientMode is enabled, reject any client commands coming from non-localhost addresses
        if LocalClientMode && isClientCommand && isNotLoopback {
            fmt.Printf("Received client command from non-localhost address: %s. Command rejected.\n", remoteAddr.IP)
            continue
        }

        if handler := commandHandlers[dg.Command]; handler != nil {
            handler(dg, remoteAddr)
        } else {
            fmt.Printf("No handler for command: %d\n", dg.Command)
        }
    }
}
