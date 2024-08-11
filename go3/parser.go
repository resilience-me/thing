package main

import (
    "bytes"        // For trimming null characters from byte slices
    "encoding/binary"
    "ripple/database"
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

// validateAndParseClientDatagram validates the datagram and returns a parsed Datagram, an error message if any, and an error object.
func validateAndParseClientDatagram(buf []byte, dg *Datagram) (string, error) {
    // Check user and peer existence
    errorMessage, err := checkUserAndPeerExist(dg)
    if err != nil {
        return errorMessage, fmt.Errorf("validation failed during user and peer existence check: %w", err)
    }

    // Load client secret key
    secretKey, err := loadClientSecretKey(dg)
    if err != nil {
        return "Error loading client secret key", fmt.Errorf("validation failed during secret key loading: %w", err)
    }

    // Verify HMAC
    if !verifyHMAC(buf, secretKey) {
        return "Error verifying HMAC", errors.New("validation failed: HMAC verification")
    }

    // Return the parsed datagram if everything is successful
    return "", nil
}

// validateAndParseServerDatagram validates the server datagram and returns a parsed Datagram and an error if any.
func validateAndParseServerDatagram(buf []byte, dg *Datagram) error {
    secretKey, err := loadServerSecretKey(dg)
    if err != nil {
        return fmt.Errorf("error loading server secret key: %w", err)
    }

    if !verifyHMAC(buf, secretKey) {
        return errors.New("error verifying HMAC for server datagram")
    }

    return nil
}

