package languages

import "github.com/charmbracelet/lipgloss"

// Go provides Go syntax highlighting
type Go struct {
	*Base
	keywords  []string
	types     []string
	constants []string
}

// NewGo creates Go highlighter
func NewGo() *Go {
	return &Go{
		Base: NewBase(),
		keywords: []string{
			"package", "import", "func", "type", "struct", "interface",
			"var", "const", "if", "else", "for", "range", "return",
			"switch", "case", "default", "break", "continue",
			"go", "defer", "select", "chan", "map",
		},
		types: []string{
			"int", "int64", "int32", "string", "bool", "float64", "float32", "byte", "rune", "error",
		},
		constants: []string{
			"true", "false", "nil",
		},
	}
}

func (g *Go) Highlight(line string) string {
	result := line

	// Highlight keywords
	for _, kw := range g.keywords {
		result = g.HighlightWord(result, kw, g.keywordStyle)
	}

	// Highlight types
	for _, t := range g.types {
		result = g.HighlightWord(result, t, g.typeStyle)
	}

	// Highlight constants
	for _, c := range g.constants {
		result = g.HighlightWord(result, c, g.constantStyle)
	}

	// Highlight strings
	result = g.HighlightStrings(result, g.stringStyle)

	// Highlight comments
	result = g.HighlightComments(result, "//", g.commentStyle)

	return lipgloss.NewStyle().Render(result)
}
