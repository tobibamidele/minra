package viewport

import "github.com/tobibamidele/minra/internal/cursor"

// AdjustScroll adjusts scroll to keep cursor visible
func (v *Viewport) AdjustScroll(cur *cursor.Cursor) {
	// Vertical scroll
	if cur.Line() < v.scrollY {
		v.scrollY = cur.Line()
	}
	if cur.Line() >= v.scrollY+v.height {
		v.scrollY = cur.Line() - v.height + 1
	}

	// Horizontal scroll
	displayCol := v.calculateDisplayCol(cur)
	if displayCol < v.scrollX {
		v.scrollX = displayCol
	}
	if displayCol >= v.scrollX+v.width {
		v.scrollX = displayCol - v.width + 1
	}
}

// ScrollUp scrolls up by lines
func (v *Viewport) ScrollUp(lines int) {
	v.scrollY -= lines
	if v.scrollY < 0 {
		v.scrollY = 0
	}
}

// ScrollDown scrolls down by lines
func (v *Viewport) ScrollDown(lines int) {
	v.scrollY += lines
	maxScroll := v.buffer.LineCount() - v.height
	if maxScroll < 0 {
		maxScroll = 0
	}
	if v.scrollY > maxScroll {
		v.scrollY = maxScroll
	}
}

// CenterCursor centers cursor in viewport
func (v *Viewport) CenterCursor(cur *cursor.Cursor) {
	v.scrollY = cur.Line() - v.height/2
	if v.scrollY < 0 {
		v.scrollY = 0
	}
}

// calculateDisplayCol calculates display column accounting for tabs
func (v *Viewport) calculateDisplayCol(cur *cursor.Cursor) int {
	line := v.buffer.Line(cur.Line())
	displayCol := 0

	for i := 0; i < cur.Col() && i < len(line); i++ {
		if line[i] == '\t' {
			displayCol += v.tabSize - (displayCol % v.tabSize)
		} else {
			displayCol++
		}
	}

	return displayCol
}
