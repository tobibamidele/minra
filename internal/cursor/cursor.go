package cursor

// import "github.com/tobibamidele/minra/internal/buffer"

// Cursor represents cursor position
type Cursor struct {
	line int
	col  int
}

// New creates a new cursor at (0,0)
func New() *Cursor {
	return &Cursor{line: 0, col: 0}
}

// Line returns current line
func (c *Cursor) Line() int {
	return c.line
}

// Col returns current column
func (c *Cursor) Col() int {
	return c.col
}

// SetPosition sets cursor position
func (c *Cursor) SetPosition(line, col int) {
	c.line = line
	c.col = col
}

// Position returns line and column
func (c *Cursor) Position() (int, int) {
	return c.line, c.col
}
