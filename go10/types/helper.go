package types

import (
    "encoding/binary"
    "syscall" // Imported to use syscall.Clen for finding the length of a null-terminated byte slice.
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

// BytesToString converts a byte slice to a string, stopping at the first null character.
func BytesToString(data []byte) string {
    length := syscall.Clen(data) // Use syscall.clen to find the length up to the first null byte
    return string(data[:length]) // Convert the trimmed byte slice to a string
}
