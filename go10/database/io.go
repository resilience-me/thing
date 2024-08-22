package database

import (
    "fmt"
    "os"
    "path/filepath"
    "strconv"
)

// ReadFile reads the content of a file and returns it as a byte slice.
func ReadFile(dir, filename string) ([]byte, error) {
    filePath := filepath.Join(dir, filename)
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
    }
    return data, nil
}

// GetUint32FromFile reads the contents of a file, parses it as a uint32, and returns the value.
func GetUint32FromFile(dir, filename string) (uint32, error) {
    data, err := ReadFile(dir, filename)
    if err != nil {
        return 0, err
    }

    value, err := strconv.ParseUint(string(data), 10, 32)
    if err != nil {
        return 0, fmt.Errorf("error parsing value from file %s: %v", filepath.Join(dir, filename), err)
    }
    return uint32(value), nil
}

// ReadTimeFromFile reads a Unix timestamp from a file and returns it as an int64.
func ReadTimeFromFile(dir, filename string) (int64, error) {
    data, err := ReadFile(dir, filename)
    if err != nil {
        return 0, err
    }

    timestamp, err := strconv.ParseInt(string(data), 10, 64)
    if err != nil {
        return 0, fmt.Errorf("error parsing timestamp from file %s: %v", filepath.Join(dir, filename), err)
    }
    return timestamp, nil
}

// WriteUint32ToFile writes a uint32 value to a file.
func WriteUint32ToFile(dir, filename string, value uint32) error {
    filePath := filepath.Join(dir, filename)
    return os.WriteFile(filePath, []byte(fmt.Sprintf("%d", value)), 0644)
}

// writeTimeToFile writes a Unix timestamp to a file.
func WriteTimeToFile(dir, filename string, timestamp int64) error {
	filePath := filepath.Join(dir, filename)
	return os.WriteFile(filePath, []byte(fmt.Sprintf("%d", timestamp)), 0644)
}
