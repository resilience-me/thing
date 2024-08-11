package main

import (
    "bytes"        // For trimming null characters from byte slices
    "encoding/binary"
    "ripple/database"
)

// checkUserAndPeerExist checks for the existence of user and peer directories.
// It returns an error message string (empty if successful) and an error object for detailed information if an error occurs.
func checkUserAndPeerExist(dg *Datagram) (string, error) {
    exists, err := database.CheckAccountExists(dg)
    if err != nil {
        return "Error checking account existence", fmt.Errorf("error checking account existence for user '%s': %v", dg.Username, err)
    } else if !exists {
        return "User account does not exist", fmt.Errorf("account directory does not exist for user '%s'", dg.Username)
    }

    exists, err = database.CheckPeerExists(dg)
    if err != nil {
        return "Error checking peer existence", fmt.Errorf("error checking peer existence for server '%s' and user '%s': %v", dg.PeerServerAddress, dg.PeerUsername, err)
    } else if !exists {
        return "Peer account does not exist", fmt.Errorf("peer directory does not exist for server '%s' and user '%s'", dg.PeerServerAddress, dg.PeerUsername)
    }

    return "", nil // No error, directories exist
}

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
