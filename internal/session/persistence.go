package session

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/tobibamidele/minra/internal/sidebar"
)

// Save saves session to a file
func SaveSession(session *Session, path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Load loads session from file
func LoadSession(path string) (*Session, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// DefaultSessionPath returns default session file path.
// This is `$HOME/.minra/sessions/{workspace}.json`
func DefaultSessionPath(workspace string) string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".minra", "sessions", filepath.Base(workspace)+".json")
}

// DefaultUIStatePath returns the default ui state
// This is `$HOME/.minra/ui/state.json`
func DefaultUIStatePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".minra", "ui", "state.json")
}

// SaveUIState saves the ui state
func SaveUIState(
	sidebar *sidebar.Sidebar,
	path string,
) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(sidebar.GetFileTreeState()), 0644)
}
