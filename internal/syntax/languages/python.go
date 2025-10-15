package languages

// Python provides Python syntax highlighting
type Python struct {
	*Base
	keywords  []string
	constants []string
}

func NewPython() *Python {
	return &Python{
		Base: NewBase(),
		keywords: []string{
			"def", "class", "if", "elif", "else", "for", "while",
			"return", "import", "from", "as", "try", "except",
			"finally", "with", "lambda", "yield", "pass", "break",
			"continue", "raise", "assert", "global", "nonlocal",
		},
		constants: []string{"True", "False", "None"},
	}
}

func (p *Python) Highlight(line string) string {
	result := line
	
	for _, kw := range p.keywords {
		result = p.HighlightWord(result, kw, p.keywordStyle)
	}
	
	for _, c := range p.constants {
		result = p.HighlightWord(result, c, p.constantStyle)
	}
	
	result = p.HighlightStrings(result, p.stringStyle)
	result = p.HighlightComments(result, "#", p.commentStyle)
	
	return result
}
