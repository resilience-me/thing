package main

import (
    "encoding/binary"
    "errors"       // For creating error messages
    "fmt"          // For formatted I/O

    "ripple/database"  // Custom package, assuming it exists in your project
)

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

func parseDatagram(buf []byte) *Datagram {
    // Assuming buf is already confirmed to be of the correct length via io.ReadFull
    datagram := &Datagram{
        Command:           buf[0],
        Username:          bytesToString(buf[1:33]),
        PeerUsername:      bytesToString(buf[33:65]),
        PeerServerAddress: bytesToString(buf[65:97]),
        Arguments:         [256]byte{},
        Counter:           BytesToUint32(buf[353:357]),  // Directly initializing the Counter
        Signature:         [32]byte{},
    }

    // Copy data into fixed-size arrays for Arguments and Signature
    copy(datagram.Arguments[:], buf[97:353])
    copy(datagram.Signature[:], buf[357:389])  // Ensure the signature is exactly 32 bytes

    return datagram
}

// checkUserAndPeerExist checks for the existence of user and peer directories
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

// validateClientDatagram validates the client datagram
func validateClientDatagram(buf []byte, dg *Datagram) (string, error) {
    errorMessage, err := checkUserAndPeerExist(dg)
    if err != nil {
        return errorMessage, fmt.Errorf("user and peer existence check failed: %w", err)
    }

    secretKey, err := loadClientSecretKey(dg)
    if err != nil {
        return "Error loading client secret key", fmt.Errorf("loading client secret key failed: %w", err)
    }

    if !verifyHMAC(buf, secretKey) {
        return "Error verifying HMAC", errors.New("HMAC verification failed")
    }

    return "", nil
}

// validateServerDatagram validates the server datagram
func validateServerDatagram(buf []byte, dg *Datagram) error {
    secretKey, err := loadServerSecretKey(dg)
    if err != nil {
        return fmt.Errorf("loading server secret key failed: %w", err)
    }

    if !verifyHMAC(buf, secretKey) {
        return errors.New("HMAC verification failed")
    }

    return nil
}
