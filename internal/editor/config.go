package editor

// Config holds app configuration
type Config struct {
	TabSize         int    // Number of space for a tab
	LineNumbers     bool   // Show line numbers
	SyntaxHighlight bool   // Highlight syntax?
	AutoSave        bool   // Auto save on doc change
	Theme           string // Theme to use
}

// DefaultConfig returns the default editor config
func DefaultConfig() *Config {
	return &Config{
		TabSize:         4,
		LineNumbers:     true,
		SyntaxHighlight: true,
		AutoSave:        false,
		Theme:           "default",
	}
}
