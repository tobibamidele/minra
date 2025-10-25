package utils

import (
	"path/filepath"
	"strings"
	"unicode"
)

var DefaultTabWidth = map[string]int{
	".go":     8,  // gofmt uses tabs, equivalent to width 8 in many editors
	".c":      8,
	".h":      8,
	".cpp":    4,
	".hpp":    4,
	".java":   4,
	".js":     2,
	".ts":     2,
	".jsx":    2,
	".tsx":    2,
	".py":     4,
	".rs":     4,
	".rb":     2,
	".php":    4,
	".html":   2,
	".css":    2,
	".scss":   2,
	".json":   2,
	".yaml":   2,
	".yml":    2,
	".xml":    2,
	".swift":  4,
	".kt":     4, // Kotlin
}

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

func GetTabSizeByFilePath(path string) int {
	ext := filepath.Ext(filepath.Base(path))
	tabSize, ok := DefaultTabWidth[ext]; if !ok {
		return 4 // Default
	}
	return tabSize
}
