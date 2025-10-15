package buffer

import (
	"strings"

	"github.com/tobibamidele/minra/internal/cursor"
)

// Buffer represents a text buffer
type Buffer struct {
	id			string
	lines		[]string
	filepath	string
	modified	bool
	cursor		*cursor.Cursor
	history		*History
	language	string
}

// New creates an empty buffer
func New() *Buffer {
	return &Buffer{
		lines: []string{""},
		filepath: "",
		modified: false,
		cursor: cursor.New(),
		history: NewHistory(),
	}
}

// NewFromContent creates a buffer from content
func NewFromContent(content string, filepath string) *Buffer {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		lines = []string{""}
	}

	return &Buffer{
		lines: lines,
		filepath: filepath,
		modified: false,
		cursor: cursor.New(),
		history: NewHistory(),
	}
}

// ID returns the buffer ID
func (b *Buffer) ID() string {
	return b.id
}

// SetID sets the buffer ID
func (b *Buffer) SetID(id string) {
	b.id = id
}

// LineCount returns the number of lines in the buffer
func (b *Buffer) LineCount() int {
	return len(b.lines)
}

// Line returns a specific line
func (b *Buffer) Line(n int) string {
	if n < 0 || n > len(b.lines) {
		return ""
	}

	return b.lines[n]
}

// Lines returns all lines
func (b *Buffer) Lines() []string {
	return b.lines
}

// SetLines sets a specific line
func (b *Buffer) SetLine(n int, content string) {
	if n >= 0 && n < len(b.lines) {
		b.lines[n] = content
		b.modified = true
	}
}

// Content returns full buffer content
func (b *Buffer) Content() string {
	return strings.Join(b.lines, "\n")
}

// Filepath returns the file path
func (b *Buffer) Filepath() string {
	return b.filepath
}

// SetFilepath sets the file path
func (b *Buffer) SetFilepath(fp string) {
	b.filepath = fp
}

// Modified reports if the buffer's been modified
func (b *Buffer) Modified() bool {
	return b.modified
}

// SetModified sets te modified flag
func (b *Buffer) SetModified(modified bool) {
	b.modified = modified
}

// Cursor returns the buffer's cursor
func (b *Buffer) Cursor() *cursor.Cursor {
	return b.cursor
}

// Language returns the detected language
func (b *Buffer) Language() string {
	return b.language
}

// SetLanguage sets the language
func (b *Buffer) SetLanguage(lang string) {
	b.language = lang
}
