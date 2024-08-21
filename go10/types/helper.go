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

// BytesToTrimmedString trims null characters from byte slices for proper string conversion.
func BytesToString(data []byte) string {
    return string(bytes.TrimRight(data, "\x00"))
}
