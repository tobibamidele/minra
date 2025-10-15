package clipboard

// Clipboard interface for copy/paste operations
type Clipboard interface {
	Copy(text string) error
	Paste() (string, error)
	Clear() error
}

// Mode represents clipboard mode
type Mode int

const (
	ModeChar Mode = iota
	ModeLine
	ModeBlock
)

// New creates appropriate clipboard (system or internal)
func New() Clipboard {
	// Try system clipboard first
	sys := NewSystem()
	if sys != nil {
		return sys
	}
	// Fallback to internal
	return NewInternal()
}
