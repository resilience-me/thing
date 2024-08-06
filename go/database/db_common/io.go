package db_common

import (
    "fmt"
    "os"
    "path/filepath"
    "strconv"
)

// GetUint32FromFile reads the contents of a file, parses it as a uint32, and returns the value.
func GetUint32FromFile(dir, filename string) (uint32, error) {
    filePath := filepath.Join(dir, filename)
    data, err := os.ReadFile(filePath)
    if err != nil {
        return 0, fmt.Errorf("error reading file %s: %v", filePath, err)
    }

    // Convert the file content to uint32
    value, err := strconv.ParseUint(string(data), 10, 32)
    if err != nil {
        return 0, fmt.Errorf("error parsing value from file %s: %v", filePath, err)
    }
    return uint32(value), nil
}

// WriteUint32ToFile writes a uint32 value to a file.
func WriteUint32ToFile(dir, filename string, value uint32) error {
    filePath := filepath.Join(dir, filename)
    return os.WriteFile(filePath, []byte(fmt.Sprintf("%d", value)), 0644)
}
