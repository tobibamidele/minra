package statusbar

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Render renders the status bar
func (s *StatusBar) Render(width int, left, right string) string {
	leftStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("240")).
		Padding(0, 1)

	rightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("240")).
		Padding(0, 1)

	leftRendered := leftStyle.Render(left)
	rightRendered := rightStyle.Render(right)

	gap := width - lipgloss.Width(leftRendered) - lipgloss.Width(rightRendered)
	if gap < 0 {
		gap = 0
	}

	statusBar := lipgloss.NewStyle().
		Background(lipgloss.Color("240")).
		Width(width).
		Render(leftRendered + strings.Repeat(" ", gap) + rightRendered)

	return statusBar
}
