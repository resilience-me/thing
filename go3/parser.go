func parseDatagram(buf []byte) (*Datagram, error) {
    // Assuming buf is already confirmed to be of the correct length via io.ReadFull
    datagram := &Datagram{
        Command:           buf[0],
        Username:          trimRightZeroes(buf[1:33]),
        PeerUsername:      trimRightZeroes(buf[33:65]),
        PeerServerAddress: trimRightZeroes(buf[65:97]),
        Arguments:         [256]byte{},
        Counter:           [4]byte{},
        Signature:         [32]byte{},
    }

    // Copy data into fixed-size arrays
    copy(datagram.Arguments[:], buf[97:353])
    copy(datagram.Counter[:], buf[353:357])
    copy(datagram.Signature[:], buf[357:389])  // Ensure the signature is exactly 32 bytes

    return datagram, nil
}

// Helper function to trim null characters from byte slices for proper string conversion
func trimRightZeroes(data []byte) string {
    return string(bytes.TrimRight(data, "\x00"))
}
