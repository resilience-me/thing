package db_server

import (
	"ripple/main"
	"ripple/database"
)

// GetCounterOut retrieves the counter_out value using the datagram to determine the directory.
func GetCounterOut(dg *main.Datagram) (uint32, error) {
	peerDir := database.GetPeerDir(dg)
	return database.GetUint32FromFile(peerDir, "counter_out.txt")
}

// GetCounterIn retrieves the counter_in value using the datagram to determine the directory.
func GetCounterIn(dg *main.Datagram) (uint32, error) {
	peerDir := database.GetPeerDir(dg)
	return database.GetUint32FromFile(peerDir, "counter_in.txt")
}

// SetCounterOut sets the counter_out value.
func SetCounterOut(dg *main.Datagram, value uint32) error {
	peerDir := database.GetPeerDir(dg)
	return database.WriteUint32ToFile(peerDir, "counter_out.txt", value)
}

// SetCounterIn sets the counter_in value.
func SetCounterIn(dg *main.Datagram) error {
	peerDir := database.GetPeerDir(dg)
	return database.WriteUint32ToFile(peerDir, "counter_in.txt", dg.Counter)
}

// ValidateCounterIn checks if the provided counter is greater than the stored counter_in value.
func ValidateCounterIn(datagram *main.Datagram) error {
    // Retrieve the stored counter_in value
    prevCounterIn, err := GetCounterIn(datagram)
    if err != nil {
        return fmt.Errorf("error getting stored counter_in for user %s: %v", datagram.Username, err)
    }

    // Check if the incoming counter is valid (greater than the stored counter_in)
    if datagram.Counter <= prevCounterIn {
        return fmt.Errorf("counter_in validation failed for user %s", datagram.Username)
    }

    return nil
}

// GetAndIncrementCounterOut retrieves the current counter_out, increments it, and updates the database.
// It returns the counter value before it was incremented.
func GetAndIncrementCounterOut(datagram *main.Datagram) (uint32, error) {
    // Retrieve the current value of counter_out from the database.
    counterOut, err := GetCounterOut(datagram)
    if err != nil {
        return 0, err  // Return error if unable to fetch the counter.
    }

    // Increment the counter and update it in the database within the same function call.
    if err := SetCounterOut(datagram, counterOut + 1); err != nil {
        return 0, err  // Return error if unable to update the counter.
    }

    // Return the original counter value that was fetched.
    return counterOut, nil
}
