package types

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
