package main

import (
    "fmt"
    "net"
    "os"
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

var commandHandlers = [256]CommandHandler{
    handleSetTrustline, // Command 0
    handleGetTrustline, // Command 1
    // other handlers explicitly set to nil by default
}

func handleSetTrustline(dg Datagram, addr *net.UDPAddr) {
    fmt.Println("Handling Set Trustline")
}

func handleGetTrustline(dg Datagram, addr *net.UDPAddr) {
    fmt.Println("Handling Get Trustline")
}

func main() {
    addr := net.UDPAddr{
        Port: 2012,
        IP:   net.ParseIP("::"),
    }
    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error listening: %v\n", err)
        return
    }
    defer conn.Close()

    fmt.Println("Server listening on all interfaces for IPv4 and IPv6.")

    for {
        var dg Datagram
        _, remoteAddr, err := conn.ReadFromUDP(dg[:])
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error reading from UDP: %v\n", err)
            continue
        }

        handler := commandHandlers[dg.Command]
        if handler != nil {
            handler(dg, remoteAddr)
        } else {
            fmt.Printf("No handler for command: %d\n", dg.Command)
        }
    }
}
