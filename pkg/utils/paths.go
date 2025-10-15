package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// GetFileName returns just the filename from a path
func GetFileName(path string) string {
	return filepath.Base(path)
}

// GetFileExtension returns the file extension
func GetFileExtension(path string) string {
	return filepath.Ext(path)
}

// GetDirectory returns the directory part of a path
func GetDirectory(path string) string {
	return filepath.Dir(path)
}

// JoinPath joins path elements
func JoinPath(elements ...string) string {
	return filepath.Join(elements...)
}

// AbsolutePath returns absolute path
func AbsolutePath(path string) (string, error) {
	return filepath.Abs(path)
}

// RelativePath returns relative path
func RelativePath(basepath, targpath string) (string, error) {
	return filepath.Rel(basepath, targpath)
}

// HomeDir returns user home directory
func HomeDir() (string, error) {
	return os.UserHomeDir()
}

// ExpandHome expands ~ to home directory
func ExpandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := HomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

// IsHidden checks if a file is hidden (starts with .)
func IsHidden(path string) bool {
	name := filepath.Base(path)
	return strings.HasPrefix(name, ".")
}
