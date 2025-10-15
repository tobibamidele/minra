package widgets

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tobibamidele/minra/internal/ui"
)

// RenameWidget represents a model input widget for renaming a file
type RenameWidget struct {
	visible bool
	input string
	cursorPos int
	originalName string
	width int
}

// Create a new RenameWIdget
func NewRenameWidget() *RenameWidget {
	return &RenameWidget{
		visible: false,
		input: "",
		cursorPos: 0,
		width: 40, // default width
	}
}

// Shows the widget with current filename prefilled
func (w *RenameWidget) Show(currentFileName string) {
	w.visible = true
	w.input = currentFileName
	w.originalName = currentFileName
	w.cursorPos = len(currentFileName)
}

// Hides the widget and clears the input
func (w *RenameWidget) Hide() {
	w.visible = false
	w.input = ""
	w.cursorPos = 0
	w.originalName = ""
}

// Returns if the current RenameWidget is visible
func (w *RenameWidget) IsVisible() bool {
	return w.visible
}

// Get input returns the current input value
func (w *RenameWidget) GetInput() string {
	return w.input
}

// This inserts a character into the current cursor position
func (w *RenameWidget) InsertRune(r rune) {
	before := w.input[:w.cursorPos]
	after := w.input[w.cursorPos:]
	w.input = before + string(r) + after
	w.cursorPos++
}

// This deletes a character at the cursor position
func (w *RenameWidget) DeleteRune() {
	if w.cursorPos > 0 { 
		before := w.input[:w.cursorPos - 1]
		after := w.input[w.cursorPos:]
		w.input = before + after
		w.cursorPos--
	}
}

// Moves the cursor left
func (w *RenameWidget) MoveCursorLeft() {
	if w.cursorPos > 0 {
		w.cursorPos--
	}
}

// Moves the cursor right
func (w *RenameWidget) MoveCursorRight() {
	if w.cursorPos < len(w.input) {
		w.cursorPos++
	}
}

// Moves the cursor to the starting position
func (w *RenameWidget) MoveCursorToStart() {
	w.cursorPos = 0
}

// Moves the cursor to the end
func (w *RenameWidget) MoveCursorToEnd() {
	w.cursorPos = len(w.input)
}

// Renders the rename widget
func (w *RenameWidget) Render() string {
	if !w.visible {
		return ""
	}

	cursorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("230")).
		Foreground(lipgloss.Color("0"))

	// Create the input field display
	displayInput := w.input

	// Insert cursor indicator
	if w.cursorPos < len(displayInput) {
		before := displayInput[:w.cursorPos]
		after := displayInput[w.cursorPos:]
		cursorChar := string(displayInput[w.cursorPos])

		lineContent := before + cursorStyle.Render(cursorChar)
		if w.cursorPos+1 < len(displayInput) {
			lineContent += after[1:]
		}

		displayInput = lineContent
	} else {
		// Cursor is at the end
		displayInput = displayInput + cursorStyle.Render(" ")
	}

	styleWidth := w.width - 4 // Account for padding and borders

	// Truncate filename is too long for widget width
	maxInputWidth := styleWidth 
	if len(displayInput) >= maxInputWidth {
		// Show the position around the cursor
		start := w.cursorPos - maxInputWidth / 2
		if start < 0 {
			start = 0
		}
		end := start + maxInputWidth
		if end > len(displayInput) {
			end = len(displayInput)
			start = end - maxInputWidth
			if start < 0 {
				start = 0
			}
		}

		displayInput = displayInput[start:end]
	}

	// Build the widget content
	var content strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
			Foreground(ui.ColorAccent).
			Bold(true).
			Align(lipgloss.Center).
			Width(styleWidth)

	content.WriteString(titleStyle.Render("Rename File"))
	content.WriteString("\n\n")

	// Input field
	inputStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("236")).
			Padding(0, 1).
			Width(styleWidth)

	content.WriteString(inputStyle.Render(displayInput))
	content.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().
			Foreground(ui.ColorComment).
			Italic(true).
			Align(lipgloss.Center).
			Width(styleWidth)

	content.WriteString(helpStyle.Render("Enter: confirm | Esc: cancel"))

	boxStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ui.ColorAccent).
			Padding(1, 2).
			Width(w.width).
			Background(lipgloss.Color("235"))

	return boxStyle.Render(content.String())
}
