package main

import (
    "encoding/binary"
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

var commandHandlers = []CommandHandler{
    handleSetTrustline, // Command 0
    handleGetTrustline, // Command 1
}

func handleSetTrustline(dg Datagram, addr *net.UDPAddr) {
    trustlineAmount := binary.LittleEndian.Uint32(dg.Arguments[:4])
    fmt.Printf("Set trustline to %d for %s\n", trustlineAmount, dg.YUsername)
}

func handleGetTrustline(dg Datagram, addr *net.UDPAddr) {
    fmt.Printf("Get trustline for %s\n", dg.YUsername)
}

func main() {
    addr := net.UDPAddr{
        Port: 2012,
        IP:   net.ParseIP("::"),  // Listen on all IPv6 addresses and dual-stack for IPv4
    }
    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer conn.Close()

    for {
        var dg Datagram
        _, remoteAddr, err := conn.ReadFromUDP(dg[:])
        if err != nil {
            fmt.Println(err)
            continue
        }

        fmt.Printf("Received datagram from %s\n", remoteAddr.String())
        dispatchCommand(dg, remoteAddr)
    }
}

func dispatchCommand(dg Datagram, addr *net.UDPAddr) {
    if int(dg.Command) < len(commandHandlers) {
        commandHandlers[dg.Command](dg, addr)
    } else {
        fmt.Printf("Unknown command: %d\n", dg.Command)
    }
}
