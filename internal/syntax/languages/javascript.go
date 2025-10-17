package languages

// JavaScript provides JavaScript syntax highlighting
type JavaScript struct {
	*Base
	keywords  []string
	constants []string
}

func NewJavaScript() *JavaScript {
	return &JavaScript{
		Base: NewBase(),
		keywords: []string{
			"function", "const", "let", "var", "if", "else", "for",
			"while", "return", "class", "extends", "import", "export",
			"from", "async", "await", "try", "catch", "finally",
			"switch", "case", "default", "break", "continue", "new",
			"this", "typeof", "instanceof",
		},
		constants: []string{"true", "false", "null", "undefined"},
	}
}

func (j *JavaScript) Highlight(line string) string {
	result := line

	for _, kw := range j.keywords {
		result = j.HighlightWord(result, kw, j.keywordStyle)
	}

	for _, c := range j.constants {
		result = j.HighlightWord(result, c, j.constantStyle)
	}

	result = j.HighlightStrings(result, j.stringStyle)
	result = j.HighlightComments(result, "//", j.commentStyle)

	return result
}

// TypeScript is an alias for JavaScript
type TypeScript struct {
	*JavaScript
}

func NewTypeScript() *TypeScript {
	return &TypeScript{
		JavaScript: NewJavaScript(),
	}
}
