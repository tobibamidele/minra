package utils

import (
	"strings"
	"unicode"
)

// TruncateString truncates a string to maxLen
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// IsWhitespace checks if a rune is whitespace
func IsWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

// StripWhitespace removes leading/trailing whitespace
func StripWhitespace(s string) string {
	return strings.TrimSpace(s)
}

// CountLines counts lines in text
func CountLines(text string) int {
	return strings.Count(text, "\n") + 1
}

// IndentLevel calculates indentation level
func IndentLevel(line string, tabSize int) int {
	level := 0
	for _, ch := range line {
		if ch == ' ' {
			level++
		} else if ch == '\t' {
			level += tabSize
		} else {
			break
		}
	}
	return level / tabSize
}

// ExpandTabs expands tabs to spaces
func ExpandTabs(line string, tabSize int) string {
	var result strings.Builder
	col := 0
	
	for _, ch := range line {
		if ch == '\t' {
			spacesToAdd := tabSize - (col % tabSize)
			for i := 0; i < spacesToAdd; i++ {
				result.WriteRune(' ')
			}
			col += spacesToAdd
		} else {
			result.WriteRune(ch)
			col++
		}
	}
	
	return result.String()
}
