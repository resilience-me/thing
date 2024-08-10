package main

import (
    "bytes"        // For trimming null characters from byte slices
    "encoding/binary"
)

func parseDatagram(buf []byte) *Datagram {
    // Assuming buf is already confirmed to be of the correct length via io.ReadFull
    datagram := &Datagram{
        Command:           buf[0],
        Username:          trimRightZeroes(buf[1:33]),
        PeerUsername:      trimRightZeroes(buf[33:65]),
        PeerServerAddress: trimRightZeroes(buf[65:97]),
        Arguments:         [256]byte{},
        Counter:           binary.BigEndian.Uint32(buf[353:357]),  // Directly initializing the Counter
        Signature:         [32]byte{},
    }

    // Copy data into fixed-size arrays for Arguments and Signature
    copy(datagram.Arguments[:], buf[97:353])
    copy(datagram.Signature[:], buf[357:389])  // Ensure the signature is exactly 32 bytes

    return datagram
}

// Helper function to trim null characters from byte slices for proper string conversion
func trimRightZeroes(data []byte) string {
    return string(bytes.TrimRight(data, "\x00"))
}
