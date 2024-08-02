package main

import (
    "fmt"
    "net"
)

type Datagram struct {
    Command        byte
    XUsername      [32]byte
    YUsername      [32]byte
    YServerAddress [32]byte
    Arguments      [256]byte
    Counter        [4]byte
    Signature      [32]byte
}

type CommandHandler func(Datagram, *net.UDPAddr)

// Direct initialization of commandHandlers
var commandHandlers = [256]CommandHandler{
    handleSetTrustline, // Command 0
    handleGetTrustline, // Command 1
    // All other handlers are implicitly set to nil
}

func main() {
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
        
        if handler := commandHandlers[dg.Command]; handler != nil {
            handler(dg, remoteAddr)
        } else {
            fmt.Printf("No handler for command: %d\n", dg.Command)
        }
    }
}

func handleSetTrustline(dg Datagram, addr *net.UDPAddr) {
    fmt.Println("Handling Set Trustline")
}

func handleGetTrustline(dg Datagram, addr *net.UDPAddr) {
    fmt.Println("Handling Get Trustline")
}
