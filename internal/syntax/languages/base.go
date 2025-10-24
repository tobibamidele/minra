package languages

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tobibamidele/minra/internal/ui"
)

var commentKeywords = []string{"TODO:", "NOTE:", "FIXME:", "BUG:", "HACK:"}

// Language interface for syntax highlighting
type Language interface {
	Highlight(line string) string
}

// Base provides common highlighting utilities
type Base struct {
	keywordStyle  lipgloss.Style
	typeStyle     lipgloss.Style
	constantStyle lipgloss.Style
	stringStyle   lipgloss.Style
	commentStyle  lipgloss.Style
	commentKeywordStyle lipgloss.Style
}

func NewBase() *Base {
	return &Base{
		keywordStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("205")),
		typeStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("117")),
		constantStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("141")),
		stringStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("150")),
		commentStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("243")),
		commentKeywordStyle: lipgloss.NewStyle().Background(ui.ColorSelection).Foreground(lipgloss.Color("#ffffff")),
	}
}

// HighlightWord highlights whole words
func (b *Base) HighlightWord(text, word string, style lipgloss.Style) string {
	inString := markStringRegions(text)

	var result strings.Builder
	i := 0
	wordLen := len(word)

	for i < len(text) {
		// Skip highlighting if inside a string
		if !inString[i] &&
			i+wordLen <= len(text) &&
			text[i:i+wordLen] == word &&
			isBoundary(text, i-1) &&
			isBoundary(text, i+wordLen) {

			result.WriteString(style.Render(word))
			i += wordLen
		} else {
			result.WriteByte(text[i])
			i++
		}
	}
	return result.String()
}

func markStringRegions(text string) []bool {
	inString := false
	quote := rune(0)
	mark := make([]bool, len(text))

	for i, ch := range text {
		if !inString && (ch == '"' || ch == '\'') {
			inString = true
			quote = ch
			mark[i] = true
		} else if inString && ch == quote && (i == 0 || text[i-1] != '\\') {
			mark[i] = true
			inString = false
		} else if inString {
			mark[i] = true
		}
	}
	return mark
}

func isBoundary(s string, idx int) bool {
	if idx < 0 || idx >= len(s) {
		return true
	}
	ch := s[idx]
	return ch != '_' && (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9')
}

// HighlightStrings highlights quoted strings
func (b *Base) HighlightStrings(text string, style lipgloss.Style) string {
	inString := false
	quote := rune(0)
	var result strings.Builder

	for i, ch := range text {
		if !inString && (ch == '"' || ch == '\'') {
			inString = true
			quote = ch
			result.WriteString(style.Render(string(ch)))
		} else if inString && ch == quote && (i == 0 || text[i-1] != '\\') {
			result.WriteString(style.Render(string(ch)))
			inString = false
		} else if inString {
			result.WriteString(style.Render(string(ch)))
		} else {
			result.WriteRune(ch)
		}
	}

	return strings.ReplaceAll(result.String(), "\t", "  ")
	// return result.String()
}

// // HighlightComments highlights comments
// func (b *Base) HighlightComments(text, commentStart string, style lipgloss.Style) string {
// 	idx := strings.Index(text, commentStart)
// 	if idx == -1 {
// 		return text
// 	}
//
// 	before := text[:idx]
// 	comment := text[idx:]
//
// 	return before + style.Render(comment)
// }

// HighlightComments highlights comments and special comment keywords
func (b *Base) HighlightComments(text, commentStart string, style lipgloss.Style) string {
	idx := strings.Index(text, commentStart)
	if idx == -1 {
		return text
	}

	before := text[:idx]
	comment := text[idx:]

	// First highlight whole comment normally
	highlighted := style.Render(comment)

	// Then highlight keywords inside the comment
	for _, kw := range commentKeywords {
		highlighted = strings.ReplaceAll(
			highlighted,
			kw,
			b.commentKeywordStyle.Render(kw),
		)
	}

	return before + highlighted
}

// TODO:
