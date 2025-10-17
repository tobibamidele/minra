package clipboard

import (
	"strings"

	"github.com/tobibamidele/minra/internal/buffer"
	"github.com/tobibamidele/minra/internal/cursor"
)

// CopySelection copies text from selection
func CopySelection(buf *buffer.Buffer, startLine, startCol, endLine, endCol int) string {
	if startLine > endLine || (startLine == endLine && startCol > endCol) {
		startLine, endLine = endLine, startLine
		startCol, endCol = endCol, startCol
	}

	if startLine == endLine {
		line := buf.Line(startLine)
		if startCol >= len(line) {
			return ""
		}
		if endCol > len(line) {
			endCol = len(line)
		}
		return line[startCol:endCol]
	}

	var result strings.Builder

	firstLine := buf.Line(startLine)
	if startCol < len(firstLine) {
		result.WriteString(firstLine[startCol:])
	}
	result.WriteString("\n")

	for i := startLine + 1; i < endLine; i++ {
		result.WriteString(buf.Line(i))
		result.WriteString("\n")
	}

	lastLine := buf.Line(endLine)
	if endCol > len(lastLine) {
		endCol = len(lastLine)
	}
	result.WriteString(lastLine[:endCol])

	return result.String()
}

// PasteAtCursor pastes text at cursor
func PasteAtCursor(buf *buffer.Buffer, cur *cursor.Cursor, text string) {
	lines := strings.Split(text, "\n")

	if len(lines) == 1 {
		for _, r := range text {
			buf.InsertRune(cur.Line(), cur.Col(), r)
			cur.MoveRight(buf)
		}
		return
	}

	for _, r := range lines[0] {
		buf.InsertRune(cur.Line(), cur.Col(), r)
		cur.MoveRight(buf)
	}

	buf.InsertNewline(cur.Line(), cur.Col())
	cur.MoveDown(buf)
	cur.MoveToLineStart()

	for i := 1; i < len(lines); i++ {
		for _, r := range lines[i] {
			buf.InsertRune(cur.Line(), cur.Col(), r)
			cur.MoveRight(buf)
		}

		if i < len(lines)-1 {
			buf.InsertNewline(cur.Line(), cur.Col())
			cur.MoveDown(buf)
			cur.MoveToLineStart()
		}
	}
}
