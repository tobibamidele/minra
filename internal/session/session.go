package session

// Session implements an editor session
type Session struct {
	Workspace	string
	OpenFiles	[]string
	ActiveFile	string
	WindowSize 	struct {
		Width	int
		Height	int
	}
}

// New creates a new session
func New(workspace string) *Session {
	return &Session{
		Workspace: workspace,
		OpenFiles: make([]string, 0),
	}
}

// AddFile adds a file to session
func (s *Session) AddFile(filepath string) {
	for _, f := range s.OpenFiles {
		if f == filepath {
			return
		}
	}

	s.OpenFiles = append(s.OpenFiles, filepath)
}

// RemoveFile removes a file from session
func (s *Session) RemoveFile(filepath string) {
	for i, f := range s.OpenFiles {
		if f == filepath {
			s.OpenFiles = append(s.OpenFiles[:i], s.OpenFiles[i+1:]...)
			return
		}
	}
}

// SetActiveFile sets the active file
func (s *Session) SetActiveFile(filepath string) {
	s.ActiveFile = filepath
}
