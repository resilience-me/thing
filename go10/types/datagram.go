package types

import "ripple/config"

// Datagram holds the structure of the incoming data
type Datagram struct {
    Command           byte
    Username          string
    PeerUsername      string
    PeerServerAddress string
    Arguments         [256]byte
    Counter           uint32
    Signature         [32]byte
}

// NewDatagram creates a new Datagram instance with the specified parameters
func NewDatagram(recipient, sender string, counter uint32) *Datagram {
    return &Datagram{
        Username:          recipient,
        PeerUsername:      sender,
        PeerServerAddress: config.GetServerAddress(),
        Counter:           counter,
    }
}
