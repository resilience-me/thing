package main

import (
    "fmt"
    "net"
    "os"
)

// Datagram structure to define the packet structure
type Datagram struct {
    Command        byte
    XUsername      [32]byte
    YUsername      [32]byte
    YServerAddress [32]byte
    Arguments      [256]byte
    Counter        [4]byte
    Signature      [32]byte
}

// Define the type for command handlers
type CommandHandler func(Datagram, *net.UDPAddr)

// Array of command handlers, assuming commands 0 and 1 are implemented
var commandHandlers = []CommandHandler{
    handleSetTrustline, // Command 0
    handleGetTrustline, // Command 1
    // More handlers can be added here sequentially
}

func main() {
    // Listening on both IPv4 and IPv6 addresses
    addr := net.UDPAddr{
        Port: 2012,
        IP:   net.ParseIP("::"), // This handles both IPv6 and IPv4
    }

    // Setup the UDP server
    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error listening on UDP port %d: %v\n", addr.Port, err)
        return
    }
    defer conn.Close()

    fmt.Printf("Server is listening on all network interfaces for both IPv4 and IPv6 at port %d\n", addr.Port)

    // Infinite loop to handle incoming datagrams
    for {
        var dg Datagram
        n, remoteAddr, err := conn.ReadFromUDP(dg[:])
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error reading from UDP: %v\n", err)
            continue
        }
        if n < len(dg) {
            fmt.Println("Received incomplete datagram")
            continue
        }

        // Dispatch the command using an array index
        if int(dg.Command) < len(commandHandlers) {
            commandHandlers[dg.Command](dg, remoteAddr)
        } else {
            fmt.Printf("Unknown command: %d\n", dg.Command)
        }
    }
}
