package cursor


type BufferReader interface {
	Line(n int) string
	LineCount() int
}


// MoveUp moves cursor up
func (c *Cursor) MoveUp(buf BufferReader) {
	if c.line > 0 {
		c.line--
		lineLen := len(buf.Line(c.line))
		if c.col > lineLen {
			c.col = lineLen
		}
	}
}

// MoveDown moves cursor down
func (c *Cursor) MoveDown(buf BufferReader) {
	if c.line < buf.LineCount()-1 {
		c.line++
		lineLen := len(buf.Line(c.line))
		if c.col > lineLen {
			c.col = lineLen
		}
	}
}

// MoveLeft moves cursor left
func (c *Cursor) MoveLeft(buf BufferReader) {
	if c.col > 0 {
		c.col--
	} else if c.line > 0 {
		c.line--
		c.col = len(buf.Line(c.line))
	}
}

// MoveRight moves cursor right
func (c *Cursor) MoveRight(buf BufferReader) {
	lineLen := len(buf.Line(c.line))
	if c.col < lineLen {
		c.col++
	} else if c.line < buf.LineCount()-1 {
		c.line++
		c.col = 0
	}
}

// MoveToLineStart moves to start of line
func (c *Cursor) MoveToLineStart() {
	c.col = 0
}

// MoveToLineEnd moves to end of line
func (c *Cursor) MoveToLineEnd(buf BufferReader) {
	c.col = len(buf.Line(c.line))
}

// MoveToBufferStart moves to start of buffer
func (c *Cursor) MoveToBufferStart() {
	c.line = 0
	c.col = 0
}

// MoveToBufferEnd moves to end of buffer
func (c *Cursor) MoveToBufferEnd(buf BufferReader) {
	c.line = buf.LineCount() - 1
	c.col = len(buf.Line(c.line))
}

// MovePageUp moves up one page
func (c *Cursor) MovePageUp(buf BufferReader, pageSize int) {
	c.line -= pageSize
	if c.line < 0 {
		c.line = 0
	}
	lineLen := len(buf.Line(c.line))
	if c.col > lineLen {
		c.col = lineLen
	}
}

// MovePageDown moves down one page
func (c *Cursor) MovePageDown(buf BufferReader, pageSize int) {
	c.line += pageSize
	if c.line >= buf.LineCount() {
		c.line = buf.LineCount() - 1
	}
	lineLen := len(buf.Line(c.line))
	if c.col > lineLen {
		c.col = lineLen
	}
}

// MoveWordForward moves to next word
func (c *Cursor) MoveWordForward(buf BufferReader) {
	line := buf.Line(c.line)
	
	// Skip current word
	for c.col < len(line) && !isWordBoundary(rune(line[c.col])) {
		c.col++
	}
	
	// Skip whitespace
	for c.col < len(line) && isWordBoundary(rune(line[c.col])) {
		c.col++
	}
	
	// If at end of line, move to next line
	if c.col >= len(line) && c.line < buf.LineCount()-1 {
		c.line++
		c.col = 0
	}
}

// MoveWordBackward moves to previous word
func (c *Cursor) MoveWordBackward(buf BufferReader) {
	if c.col == 0 {
		if c.line > 0 {
			c.line--
			c.col = len(buf.Line(c.line))
		}
		return
	}
	
	line := buf.Line(c.line)
	c.col--
	
	// Skip whitespace
	for c.col > 0 && isWordBoundary(rune(line[c.col])) {
		c.col--
	}
	
	// Skip word
	for c.col > 0 && !isWordBoundary(rune(line[c.col-1])) {
		c.col--
	}
}

func isWordBoundary(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '.' || r == ',' || r == ';' || r == ':' || r == '(' || r == ')' || r == '[' || r == ']' || r == '{' || r == '}'
}
