package widgets

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tobibamidele/minra/internal/ui"
)

// CreateWidget handles file creation
type CreateWidget struct {
	visible   bool
	input     string
	cursorPos int
	width     int
}

// NewCreateWidget creates a new create widget
func NewCreateWidget() *CreateWidget {
	return &CreateWidget{
		visible: false,
		width:   40,
	}
}

func (w *CreateWidget) Show() {
	w.visible = true
	w.input = ""
	w.cursorPos = 0
}

func (w *CreateWidget) Hide() {
	w.visible = false
	w.input = ""
	w.cursorPos = 0
}

func (w *CreateWidget) IsVisible() bool {
	return w.visible
}

func (w *CreateWidget) GetInput() string {
	return w.input
}

func (w *CreateWidget) InsertRune(r rune) {
	before := w.input[:w.cursorPos]
	after := w.input[w.cursorPos:]
	w.input = before + string(r) + after
	w.cursorPos++
}

func (w *CreateWidget) DeleteRune() {
	if w.cursorPos > 0 {
		before := w.input[:w.cursorPos-1]
		after := w.input[w.cursorPos:]
		w.input = before + after
		w.cursorPos--
	}
}

func (w *CreateWidget) MoveCursorLeft() {
	if w.cursorPos > 0 {
		w.cursorPos--
	}
}

func (w *CreateWidget) MoveCursorRight() {
	if w.cursorPos < len(w.input) {
		w.cursorPos++
	}
}

func (w *CreateWidget) MoveCursorToStart() {
	w.cursorPos = 0
}

func (w *CreateWidget) MoveCursorToEnd() {
	w.cursorPos = len(w.input)
}

func (w *CreateWidget) Render() string {
	if !w.visible {
		return ""
	}
	
	var content strings.Builder
	styleWidth := w.width - 4
	
	titleStyle := lipgloss.NewStyle().
		Foreground(ui.ColorAccent).
		Bold(true).
		Align(lipgloss.Center).
		Width(styleWidth)
	
	content.WriteString(titleStyle.Render("Create File"))
	content.WriteString("\n\n")
	
	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("236")).
		Padding(0, 1).
		Width(styleWidth)
	
	content.WriteString(inputStyle.Render(w.input))
	content.WriteString("\n\n")
	
	helpStyle := lipgloss.NewStyle().
		Foreground(ui.ColorComment).
		Italic(true).
		Align(lipgloss.Center).
		Width(styleWidth)
	
	content.WriteString(helpStyle.Render("Enter: create | Esc: cancel"))
	
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorAccent).
		Padding(1, 2).
		Width(w.width).
		Background(ui.ColorBackground)
	
	return boxStyle.Render(content.String())
}
