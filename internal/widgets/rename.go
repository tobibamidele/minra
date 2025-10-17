package widgets

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tobibamidele/minra/internal/ui"
)

// RenameWidget represents a model input widget for renaming a file
type RenameWidget struct {
	visible      bool
	input        string
	cursorPos    int
	scrollStart  int // controls horizontal scrolling offset
	originalName string
	width        int
}

// NewRenameWidget creates a new RenameWidget
func NewRenameWidget() *RenameWidget {
	return &RenameWidget{
		visible:     false,
		input:       "",
		cursorPos:   0,
		scrollStart: 0,
		width:       40, // default width
	}
}

// Show shows the widget with the current filename prefilled
func (w *RenameWidget) Show(currentFileName string) {
	w.visible = true
	w.input = currentFileName
	w.originalName = currentFileName
	w.cursorPos = len(currentFileName)
	w.scrollStart = 0
}

// Hide hides the widget and clears the input
func (w *RenameWidget) Hide() {
	w.visible = false
	w.input = ""
	w.cursorPos = 0
	w.scrollStart = 0
	w.originalName = ""
}

// IsVisible returns whether the RenameWidget is visible
func (w *RenameWidget) IsVisible() bool {
	return w.visible
}

// GetInput returns the current input value
func (w *RenameWidget) GetInput() string {
	return w.input
}

// InsertRune inserts a character at the cursor position
func (w *RenameWidget) InsertRune(r rune) {
	before := w.input[:w.cursorPos]
	after := w.input[w.cursorPos:]
	w.input = before + string(r) + after
	w.cursorPos++
}

// DeleteRune deletes the character before the cursor position
func (w *RenameWidget) DeleteRune() {
	if w.cursorPos > 0 {
		before := w.input[:w.cursorPos-1]
		after := w.input[w.cursorPos:]
		w.input = before + after
		w.cursorPos--
	}
}

// MoveCursorLeft moves the cursor one position to the left
func (w *RenameWidget) MoveCursorLeft() {
	if w.cursorPos > 0 {
		w.cursorPos--
	}
}

// MoveCursorRight moves the cursor one position to the right
func (w *RenameWidget) MoveCursorRight() {
	if w.cursorPos < len(w.input) {
		w.cursorPos++
	}
}

// MoveCursorToStart moves the cursor to the beginning
func (w *RenameWidget) MoveCursorToStart() {
	w.cursorPos = 0
}

// MoveCursorToEnd moves the cursor to the end
func (w *RenameWidget) MoveCursorToEnd() {
	w.cursorPos = len(w.input)
}

// Render renders the rename widget
func (w *RenameWidget) Render(width int) string {
	if !w.visible {
		return ""
	}

	cursorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("230")).
		Foreground(lipgloss.Color("0"))

	styleWidth := width - 10 // Account for borders/padding
	maxInputWidth := styleWidth

	// Ensure scroll follows cursor position
	if w.cursorPos < w.scrollStart {
		w.scrollStart = w.cursorPos
	} else if w.cursorPos > w.scrollStart+maxInputWidth-1 {
		w.scrollStart = w.cursorPos - maxInputWidth + 1
	}

	// Clamp scrollStart to valid range
	if w.scrollStart < 0 {
		w.scrollStart = 0
	}
	if w.scrollStart > len(w.input) {
		w.scrollStart = len(w.input)
	}

	// Compute visible slice of input
	end := w.scrollStart + maxInputWidth
	if end > len(w.input) {
		end = len(w.input)
	}
	visibleInput := w.input[w.scrollStart:end]

	// Apply cursor styling within visible window
	cursorVisiblePos := w.cursorPos - w.scrollStart
	if cursorVisiblePos >= 0 && cursorVisiblePos < len(visibleInput) {
		before := visibleInput[:cursorVisiblePos]
		after := visibleInput[cursorVisiblePos:]
		cursorChar := string(visibleInput[cursorVisiblePos])
		visibleInput = before + cursorStyle.Render(cursorChar)
		if len(after) > 1 {
			visibleInput += after[1:]
		}
	} else if cursorVisiblePos == len(visibleInput) {
		visibleInput = visibleInput + cursorStyle.Render(" ")
	}

	// --- Build the visual box ---
	var content strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(ui.ColorAccent).
		Bold(true).
		Align(lipgloss.Center).
		Width(styleWidth)

	// Outer box style
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorAccent).
		Width(width - 2).
		Background(lipgloss.Color("235"))

	// Input field
	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("236")).
		Padding(0, 1).
		Width(styleWidth)

	content.WriteString(titleStyle.Render("Rename File"))
	content.WriteString("\n")
	content.WriteString(inputStyle.Render(visibleInput))

	return boxStyle.Render(content.String())
}
