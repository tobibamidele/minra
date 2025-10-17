package fileio

import (
	"os"
	"path/filepath"
)

// WriteFile writes content to a file
func WriteFile(filepath string, content string) error {
	// Ensure directory exists
	dir := filepath[:len(filepath)-len(filepath[len(filepath)-len(filepath):])]
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return os.WriteFile(filepath, []byte(content), 0644)
}

// RenameFile renames a file
func RenameFile(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// DeleteFile deletes a file
func DeleteFile(filepath string) error {
	return os.Remove(filepath)
}

// CreateFile creates a new empty file
func CreateFile(fp string) error {
	dir := filepath.Dir(fp)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}
