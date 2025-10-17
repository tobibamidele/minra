package widgets

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tobibamidele/minra/internal/ui"
)

// SearchWidget handles search
type SearchWidget struct {
	visible   bool
	input     string
	cursorPos int
	width     int
}

// NewSearchWidget creates a new search widget
func NewSearchWidget() *SearchWidget {
	return &SearchWidget{
		visible: false,
		width:   50,
	}
}

func (w *SearchWidget) Show() {
	w.visible = true
	w.input = ""
	w.cursorPos = 0
}

func (w *SearchWidget) Hide() {
	w.visible = false
	w.input = ""
	w.cursorPos = 0
}

func (w *SearchWidget) IsVisible() bool {
	return w.visible
}

func (w *SearchWidget) GetInput() string {
	return w.input
}

func (w *SearchWidget) InsertRune(r rune) {
	before := w.input[:w.cursorPos]
	after := w.input[w.cursorPos:]
	w.input = before + string(r) + after
	w.cursorPos++
}

func (w *SearchWidget) DeleteRune() {
	if w.cursorPos > 0 {
		before := w.input[:w.cursorPos-1]
		after := w.input[w.cursorPos:]
		w.input = before + after
		w.cursorPos--
	}
}

func (w *SearchWidget) Render() string {
	if !w.visible {
		return ""
	}

	var content strings.Builder
	styleWidth := w.width - 4

	titleStyle := lipgloss.NewStyle().
		Foreground(ui.ColorWarning).
		Bold(true).
		Align(lipgloss.Center).
		Width(styleWidth)

	content.WriteString(titleStyle.Render("Search"))
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

	content.WriteString(helpStyle.Render("Enter: search | Esc: cancel"))

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.ColorWarning).
		Padding(1, 2).
		Width(w.width).
		Background(lipgloss.Color("235"))

	box := boxStyle.Render(content.String())

	return box
}
