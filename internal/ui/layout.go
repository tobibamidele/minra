package ui

// Layout calculates layout dimensions
type Layout struct {
	Width  int
	Height int
}

// New creates a new layout
func NewLayout(width, height int) *Layout {
	return &Layout{
		Width:  width,
		Height: height,
	}
}

// SidebarWidth calculates sidebar width
func (l *Layout) SidebarWidth() int {
	return 30
}

// ViewportWidth calculates viewport width
func (l *Layout) ViewportWidth(sidebarVisible bool) int {
	width := l.Width
	if sidebarVisible {
		width -= l.SidebarWidth()
	}
	return width - 6 // Borders and padding
}

// ViewportHeight calculates viewport height
func (l *Layout) ViewportHeight() int {
	return l.Height - 4 // Tab bar, status bar, borders
}

// TabBarHeight returns tab bar height
func (l *Layout) TabBarHeight() int {
	return 1
}

// StatusBarHeight returns status bar height
func (l *Layout) StatusBarHeight() int {
	return 1
}
