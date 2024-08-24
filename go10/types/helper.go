package types

import (
    "encoding/binary"
)

// Based on syscall.Clen for finding the length of a null-terminated byte slice.
func clen(n []byte) int {
    for i := 0; i < len(n); i++ {
        if n[i] == 0 {
            return i
        }
    }
    return len(n)
}

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

// BytesToString converts a byte slice to a string, stopping at the first null character.
func BytesToString(data []byte) string {
    length := clen(data) // Use clen to find the length up to the first null byte
    return string(data[:length]) // Convert the trimmed byte slice to a string
}

// PadStringTo32Bytes pads a string into a 32-byte byte slice.
func PadStringTo32Bytes(str string) []byte {
    // Create a byte slice of size 32, initialized to zero values
    paddedSlice := make([]byte, 32)
    // Copy the contents of the string into the slice
    copy(paddedSlice, str)
    return paddedSlice
}
