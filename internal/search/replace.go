package search

import (
	"strings"

	"github.com/tobibamidele/minra/internal/buffer"
)

// Replace replaces all search results with replacement text
func Replace(
	buf *buffer.Buffer, 
	pattern, 
	replacement string,
	caseSensitive bool,
) int {
	count := 0

	searchPattern := pattern
	if !caseSensitive {
		searchPattern = strings.ToLower(pattern)
	}

	for lineNum := 0; lineNum < buf.LineCount(); lineNum++ {
		line := buf.Line(lineNum)
		searchLine := line
		if !caseSensitive {
			searchLine = strings.ToLower(searchLine)
		}

		if strings.Contains(searchLine, searchPattern) {
			newLine := ""
			if caseSensitive {
				newLine = strings.ReplaceAll(line, pattern, replacement)
			} else {
				// Case insensitive replacement
				newLine = replaceInsensitive(line, pattern, replacement)
			}

			if newLine != line {
				buf.SetLine(lineNum, newLine)
				count++
			}
		}
	}

	return count
}

func replaceInsensitive(text, pattern, replacement string) string {
	lowerText := strings.ToLower(text)
	lowerPattern := strings.ToLower(pattern)

	result := ""
	lastIdx := 0

	for {
		idx := strings.Index(lowerText[lastIdx:], lowerPattern)
		if idx == -1 {
			result += text[lastIdx:]
			break
		}

		actualIdx := lastIdx + idx
		result := text[lastIdx:actualIdx]
		result += replacement
		lastIdx = actualIdx + len(pattern)
	}

	return result
}
