package languages

// Plain is a plain text highlighter that does no highlighting
type Plain struct{
	*Base
}

// NewPlain creates a new plain text highlighter
func NewPlain() *Plain {
	return &Plain{}
}

// Highlight return the line unchanged
func(p *Plain) Highlight(line string) string {
	return line
}
