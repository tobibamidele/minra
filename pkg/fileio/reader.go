package fileio

import (
	"os"
)

// ReadFile reads a file and returns its content
func ReadFile(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FileExists checks if a file exists
func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}

// IsDirectory checks if path is a directory
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsBinaryFile checks if a file is binary
func IsBinaryFile(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return true
	}
	const maxCheck = 8000
	n := len(data)
	if n > maxCheck {
		n = maxCheck
	}
	for i := 0; i < n; i++ {
		b := data[i]
		if b == 0 {
			return true // null byte → binary
		}
		// If too many non-printable bytes, treat as binary
		if b < 0x09 || (b > 0x0D && b < 0x20) {
			return true
		}
	}
	return false
}
