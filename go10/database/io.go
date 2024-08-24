package database

import (
    "fmt"
    "os"
    "io/ioutil"
    "path/filepath"
    "strconv"
)

// ReadFile reads the content of a file and returns it as a byte slice.
func ReadFile(dir, filename string) ([]byte, error) {
    filePath := filepath.Join(dir, filename)
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
    }
    return data, nil
}

// WriteFile writes a byte slice to a file.
func WriteFile(dir, filename string, data []byte) error {
    filePath := filepath.Join(dir, filename)
    return ioutil.WriteFile(filePath, data, 0644)
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
    return WriteFile(dir, filename, []byte(fmt.Sprintf("%d", value)))
}

// WriteTimeToFile writes a Unix timestamp to a file.
func WriteTimeToFile(dir, filename string, timestamp int64) error {
    return WriteFile(dir, filename, []byte(fmt.Sprintf("%d", timestamp)))
}
