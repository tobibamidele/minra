package buffer

import "strings"

// InsertRune inserts a rune at cursor position
func (b *Buffer) InsertRune(line, col int, r rune) {
	if line < 0 || line >= len(b.lines) {
		return
	}

	currentLine := b.lines[line]
	if col < 0 {
		col = 0
	}
	if col > len(currentLine){
		col = len(currentLine)
	}

	b.lines[line] = currentLine[:col] + string(r) + currentLine[col:]
	b.modified = true
}

// DeleteRune deletes a rune at a position (backspace)
func (b *Buffer) DeleteRune(line, col int) {
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

// InsertNewline inserts a newline at position
func (b *Buffer) InsertNewline(line, col int) {
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

	b.lines[line] = leftPart

	newLines := make([]string, len(b.lines)+1)
	copy(newLines[:line+1], b.lines[:line+1])
	newLines[line+1] = rightPart
	copy(newLines[line+2:], b.lines[line+1:])
	b.lines = newLines

	b.modified = true
}

// DeleteLine deletes an entire line
func (b *Buffer) DeleteLine(line int) {
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
