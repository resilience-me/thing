package handlers

import (
    "ripple/types"
    "ripple/auth"
    "ripple/types"

)

// PrepareDatagram prepares common Datagram fields and increments counter_out.
func PrepareDatagram(datagram *types.Datagram) (*types.Datagram, error) {
    // Retrieve and increment the counter_out value
    counterOut, err := db_server.GetAndIncrementCounterOut(datagram)
    if err != nil {
        return nil, fmt.Errorf("error handling counter_out for user %s: %v", datagram.Username, err)
    }

    dg := types.NewDatagram(datagram.Username, counterOut)

    return dg, nil
}

// PrepareDatagramWithRecipient prepares datagram with recipient
func PrepareDatagramWithRecipient(datagram *types.Datagram) (*types.Datagram, error) {
    // Prepare the datagram
    dgOut, err := handlers.PrepareDatagram(datagram)
    if err != nil {
        return nil, err
    }
    dgOut.Username = datagram.PeerUsername

    return dgOut, nil
}

// SignDatagram creates a signed datagram by serializing it and adding a signature.
// It requires the session to load the secret key for HMAC generation.
func SignDatagram(session main.Session, dg *types.Datagram) ([]byte, error) {
    // Serialize the datagram without the signature field
    serializedData, err := types.SerializeDatagram(dg)
    if err != nil {
        return nil, fmt.Errorf("failed to serialize datagram: %w", err)
    }

    // Load the secret key for HMAC generation
    secretKey, err := auth.LoadServerSecretKey(session.Datagram)
    if err != nil {
        return nil, fmt.Errorf("failed to load server secret key: %w", err)
    }

    // Generate HMAC for the serialized data
    signature, err := auth.GenerateHMAC(serializedData, secretKey)
    if err != nil {
        return nil, fmt.Errorf("failed to generate HMAC: %w", err)
    }

    // Update the datagram's signature field with the generated signature
    copy(dg.Signature[:], []byte(signature)) // Ensure we copy the signature into the byte array

    // Return the serialized data including the signature
    return serializedData, nil
}
