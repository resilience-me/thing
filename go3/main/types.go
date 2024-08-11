package main

import (
    "encoding/binary"
)

// uint32ToBytes converts a uint32 value to a byte slice.
func uint32ToBytes(value uint32) []byte {
    responseData := make([]byte, 4) // Allocate 4 bytes for the uint32 value
    binary.BigEndian.PutUint32(responseData, value) // Convert the uint32 to bytes
    return responseData
}

// Helper function to trim null characters from byte slices for proper string conversion
func trimRightZeroes(data []byte) string {
    return string(bytes.TrimRight(data, "\x00"))
}
