package viewport

import (
	"os"

	"github.com/tobibamidele/minra/internal/buffer"
	"golang.org/x/term"
)

// Viewport manages the visible area
type Viewport struct {
	buffer      *buffer.Buffer
	isBinary	bool
	width       int
	height      int
	scrollX     int
	scrollY     int
	lineNumbers bool
	tabSize     int
}

// New creates a new viewport
func New(buf *buffer.Buffer, width, height int) *Viewport {
	return &Viewport{
		buffer:      buf,
		width:       width,
		height:      height,
		scrollX:     0,
		scrollY:     0,
		lineNumbers: true,
		tabSize:     4,
	}
}

// SetBuffer sets the buffer
func (v *Viewport) SetBuffer(buf *buffer.Buffer) {
	v.buffer = buf
	v.scrollX = 0
	v.scrollY = 0
}

// SetSize sets viewport size
func (v *Viewport) SetSize(width, height int) {
	v.width = width
	v.height = height
}

// Width returns viewport width
func (v *Viewport) Width() int {
	return v.width
}

// Height returns viewport height
func (v *Viewport) Height() int {
	return v.height
}

// ScrollX returns horizontal scroll offset
func (v *Viewport) ScrollX() int {
	return v.scrollX
}

// ScrollY returns vertical scroll offset
func (v *Viewport) ScrollY() int {
	return v.scrollY
}

// SetTabSize sets tab size
func (v *Viewport) SetTabSize(size int) {
	v.tabSize = size
}

// TabSize returns tab size
func (v *Viewport) TabSize() int {
	return v.tabSize
}

// ToggleLineNumbers toggles line numbers
func (v *Viewport) ToggleLineNumbers() {
	v.lineNumbers = !v.lineNumbers
}

// LineNumbers returns if line numbers are shown
func (v *Viewport) LineNumbers() bool {
	return v.lineNumbers
}

// IsBinary returns if the current buffer is a binary file
func (v *Viewport) IsBinary() bool {
	return v.isBinary
}

// SetIsBinary controls whether or not the content on the buffer is binary
func (v *Viewport) SetIsBinary(value bool) {
	v.isBinary = value
}

// Returns the current terminal width in columns
func ScreenWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width == 0 {
		// Fall back if size is not detected
		return 80
	}

	return width
}

// ScreenHeight returns the current terminal height in rows
func ScreenHeight() int {
	_, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || height == 0 {
		// fallback if unable to detect size
		return 24
	}

	return height
}
