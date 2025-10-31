package buffer

import (
	"path/filepath"
	"strings"
)

var autoPairMap = map[rune]rune{
	'{':  '}',
	'[':  ']',
	'(':  ')',
	'\'': '\'',
	'"':  '"',
	'`':  '`',
}

// Contains the extension of files that should have the tags automatically closed
var autoPairTagExt = []string{".jsx", ".tsx", ".xml", ".html"}

// InsertRune inserts a rune at cursor position
func (b *Buffer) InsertRune(line, col int, r rune) {
	if !b.IsModifiable() { return }
	if line < 0 || line >= len(b.lines) {
		return
	}

	currentLine := b.lines[line]
	if col < 0 {
		col = 0
	}
	if col > len(currentLine) {
		col = len(currentLine)
	}

	// Handle tag auto-close: <tag> → </tag>
	if r == '>' {
		for i := range len(autoPairTagExt) {
			if autoPairTagExt[i] == filepath.Ext(b.Filepath()) {
				insert := b.autoCloseTags(currentLine, col)
				b.lines[line] = currentLine[:col] + insert + currentLine[col:]
				b.modified = true
				return
			}
		}
	}

	// Regular auto-pairing for brackets and quotes
	if closing, ok := autoPairMap[r]; ok {
		b.lines[line] = currentLine[:col] + string(r) + string(closing) + currentLine[col:]
		b.modified = true
		return
	}

	// Normal character insertion
	b.lines[line] = currentLine[:col] + string(r) + currentLine[col:]
	b.modified = true
}

// DeleteRune deletes a rune at a position (backspace)
func (b *Buffer) DeleteRune(line, col int) {
	if !b.IsModifiable() { return }
	if line < 0 || line >= len(b.lines) {
		return
	}

	currentLine := b.lines[line]

	// Merge with previous line if at start
	if col == 0 {
		if line > 0 {
			prevLine := b.lines[line-1]
			b.lines[line-1] = prevLine + currentLine
			b.lines = append(b.lines[:line], b.lines[line+1:]...)
			b.modified = true
		}
		return
	}

	// Delete the character before cursor
	if col > 0 && col <= len(currentLine) {
		b.lines[line] = currentLine[:col-1] + currentLine[col:]
		b.modified = true
	}
}

// InsertNewline inserts a newline at position with auto-indentation
func (b *Buffer) InsertNewline(line, col int) {
	if !b.IsModifiable() { return }
	if line < 0 || line >= len(b.lines) {
		return
	}

	currentLine := b.lines[line]
	if col < 0 {
		col = 0
	}
	if col > len(currentLine) {
		col = len(currentLine)
	}

	leftPart := currentLine[:col]
	rightPart := currentLine[col:]


	// --- Detect current indentation ---
	currentIndent := countLeadingTabsOrSpaces(leftPart)

	// --- Determine if we should increase indent (after '{', '[', '(') ---
	trimmedLeft := strings.TrimSpace(leftPart)
	shouldIncrease := strings.HasSuffix(trimmedLeft, "{") ||
		strings.HasSuffix(trimmedLeft, "[") ||
		strings.HasSuffix(trimmedLeft, "(") || 
		strings.HasSuffix(trimmedLeft, ":")

	// --- Create the indentation strings ---
	baseIndent := makeIndent(currentIndent)
	increasedIndent := makeIndent(currentIndent + b.indentWidth())

	if shouldIncrease {
		// Split line into two
		b.lines[line] = leftPart

		// Construct the new lines properly
		newLines := make([]string, 0, len(b.lines)+2)
		newLines = append(newLines, b.lines[:line+1]...)
		newLines = append(newLines, increasedIndent)      // middle indented line
		newLines = append(newLines, baseIndent+rightPart) // next line with closing brace or continuation
		newLines = append(newLines, b.lines[line+1:]...)  // rest of document

		b.lines = newLines
		b.cursor.SetPosition(line, len(increasedIndent))
	} else {
		// Normal newline
		b.lines[line] = leftPart
		newLines := make([]string, 0, len(b.lines)+1)
		newLines = append(newLines, b.lines[:line+1]...)
		newLines = append(newLines, baseIndent+rightPart)
		newLines = append(newLines, b.lines[line+1:]...)
		b.lines = newLines
	}

	b.modified = true
}

// makeIndent builds a string of tabs/spaces matching indentation width
func makeIndent(width int) string {
	return strings.Repeat(" ", width)
}

// DeleteLine deletes an entire line
func (b *Buffer) DeleteLine(line int) {
	if !b.IsModifiable() { return }
	if line < 0 || line >= len(b.lines) {
		return
	}

	b.lines = append(b.lines[:line], b.lines[line+1:]...)
	if len(b.lines) == 0 {
		b.lines = []string{""}
	}
	b.modified = true
}

// InsertText inserts text at position
func (b *Buffer) InsertText(line, col int, text string) {
	if !b.IsModifiable() { return }
	lines := strings.Split(text, "\n")

	if len(lines) == 1 {
		// Single line insert
		for _, r := range text {
			b.InsertRune(line, col, r)
			col++
		}
		return
	}

	// Multi-line insert
	currentLine := b.Line(line)
	before := currentLine[:col]
	after := currentLine[col:]

	b.SetLine(line, before+lines[0])

	for i := 1; i < len(lines)-1; i++ {
		b.InsertNewline(line+i-1, len(b.Line(line+i-1)))
		b.SetLine(line+i, lines[i])
	}

	b.InsertNewline(line+len(lines)-2, len(b.Line(line+len(lines)-2)))
	b.SetLine(line+len(lines)-1, lines[len(lines)-1]+after)
}

// countLeadingTabsOrSpaces counts indentation width
func countLeadingTabsOrSpaces(s string) int {
	count := 0
	for _, r := range s {
		if r == '\t' {
			count += 4 // treat tab as 4 spaces for alignment
		} else if r == ' ' {
			count++
		} else {
			break
		}
	}
	return count
}

// indentWidth returns your desired indent width
func (b *Buffer) indentWidth() int {
	if b.tabSize > 0 {
		return b.tabSize
	}
	return 4 // default 4 spaces if unspecified
}

// autoCloseTags creates the closing tag for html and xml keys
func (b *Buffer) autoCloseTags(currentLine string, col int) string {
	openStart := strings.LastIndex(currentLine[:col], "<")
	if openStart != -1 && openStart < col {
		tagContent := currentLine[openStart+1 : col]
		// Extract tags (e.g. <div>, <p id="x"> -> div)
		tagName := ""
		for _, ch := range tagContent {
			if ch == ' ' || ch == '\t' || ch == '>' || ch == '/' {
				break
			}
			tagName += string(ch)
		}
		if tagName != "" && !strings.HasPrefix(tagName, "/") {
			// Insert ">" + "</tagName>"
			insert := ">" + "</" + tagName + ">"
			return insert
		}
	}

	return ">"
}
