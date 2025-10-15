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
