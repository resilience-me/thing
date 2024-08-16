package db_client

import (
    "fmt"
    "ripple/main"
    "ripple/database"
)

// GetCounter retrieves the counter value using the datagram to determine the directory.
func GetCounter(dg *main.Datagram) (uint32, error) {
	accountDir := database.GetAccountDir(dg)
	return database.GetUint32FromFile(accountDir, "counter.txt")
}

// SetCounter sets the counter value.
func SetCounter(dg *main.Datagram) error {
	accountDir := database.GetAccountDir(dg)
	return database.WriteUint32ToFile(accountDir, "counter.txt", datagram.Counter)
}

// ValidateCounter checks if the provided counter is greater than the stored counter value for counter.
func ValidateCounter(datagram *main.Datagram) error {
    // Retrieve the stored counter value
    prevCounter, err := GetCounter(datagram)
    if err != nil {
        return fmt.Errorf("error getting stored counter for user %s: %v", datagram.Username, err)
    }

    // Check if the incoming counter is valid (greater than the stored counter)
    if datagram.Counter <= prevCounter {
        return fmt.Errorf("counter validation failed for user %s", datagram.Username)
    }

    return nil
}
