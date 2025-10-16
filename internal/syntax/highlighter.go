package syntax

import (
	"github.com/tobibamidele/minra/internal/syntax/languages"
)

// Highlighter provides syntax highlighting
type Highlighter struct {
	language languages.Language
}

// New creates a new highlighter
func New() *Highlighter {
	return &Highlighter{
		language: languages.NewPlain(),
	}
}

// ForExtension returns highlighter for file extension
func (h *Highlighter) ForExtension(ext string) *Highlighter {
	switch ext {
	case ".go":
		h.language = languages.NewGo()
	case ".py":
		h.language = languages.NewPython()
	case ".js", ".jsx":
		h.language = languages.NewJavaScript()
	case ".ts", ".tsx":
		h.language = languages.NewTypeScript()
	default:
		h.language = languages.NewPlain()
	}
	return h
}

// Highlight applies syntax highlighting to a line
func (h *Highlighter) Highlight(line string) string {
	if h.language == nil {
		return line
	}
	return h.language.Highlight(line)
}
