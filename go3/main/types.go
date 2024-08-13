package main

import (
    "bytes"        // For trimming null characters from byte slices
    "encoding/binary"
)

// uint32ToBytes converts a uint32 value to a byte slice.
func Uint32ToBytes(value uint32) []byte {
    responseData := make([]byte, 4) // Allocate 4 bytes for the uint32 value
    binary.BigEndian.PutUint32(responseData, value) // Convert the uint32 to bytes
    return responseData
}

// BytesToUint32 converts a byte slice to a uint32 value.
// It assumes the byte slice has at least 4 bytes.
func BytesToUint32(data []byte) uint32 {
    return binary.BigEndian.Uint32(data[:4])
}

// bytesToTrimmedString trims null characters from byte slices for proper string conversion.
func bytesToString(data []byte) string {
    return string(bytes.TrimRight(data, "\x00"))
}

// Pad to 32 bytes and copy the string's bytes into the new slice.
func PadTo32Bytes(str string) []byte {
    // Create a slice of 32 bytes initialized with zero.
    result := make([]byte, 32)
    // Convert string to bytes and copy into the result slice.
    copy(result, []byte(str))
    return result
}
