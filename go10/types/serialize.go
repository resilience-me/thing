package types

import "encoding/binary"

// SerializeDatagram converts a Datagram struct to a byte slice.
func SerializeDatagram(dg *Datagram) ([]byte, error) {
    // Create the byte slice
    data := make([]byte, 389)
    data[0] = dg.Command // First byte is the Command

    // Copy Usernames and Server Address
    copy(data[1:], dg.Username)
    copy(data[33:], dg.PeerUsername)
    copy(data[65:], dg.PeerServerAddress)

    // Write the Counter
    binary.BigEndian.PutUint32(data[353:], dg.Counter)

    return data, nil
}

func DeserializeDatagram(buf []byte) *Datagram {
    // Assuming buf is already confirmed to be of the correct length
    datagram := &Datagram{
        Command:           buf[0],
        Username:          BytesToString(buf[1:33]),
        PeerUsername:      BytesToString(buf[33:65]),
        PeerServerAddress: BytesToString(buf[65:97]),
        Arguments:         [256]byte{},
        Counter:           binary.BigEndian.Uint32(buf[353:357]),
        Signature:         [32]byte{},
    }

    // Copy data into fixed-size arrays for Arguments and Signature
    copy(datagram.Arguments[:], buf[97:353])
    copy(datagram.Signature[:], buf[357:389])

    return datagram
}
