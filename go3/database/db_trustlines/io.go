package db_trustlines

import (
    "ripple/main"
)

// GetAndIncrementCounterOut retrieves the current counter_out, increments it, and updates the database.
func GetAndIncrementCounterOut(datagram *main.Datagram) (uint32, error) {
    counterOut, err := GetCounterOut(datagram)
    if err != nil {
        return 0, err
    }

    // Increment and update the counter_out in the database
    if err := SetCounterOut(datagram, counterOut+1); err != nil {
        return 0, err
    }

    return counterOut, nil
}
