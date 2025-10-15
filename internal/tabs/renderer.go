package tabs

import (
	"strings"
	
	"github.com/charmbracelet/lipgloss"
)

// Render renders the tab bar
func (m *Manager) Render(width int) string {
	if len(m.tabs) == 0 {
		return ""
	}
	
	var b strings.Builder
	
	activeStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 2).
		Bold(true)
	
	inactiveStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("237")).
		Foreground(lipgloss.Color("250")).
		Padding(0, 2)
	
	modifiedIndicator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Render("[+]")
	
	usedWidth := 0
	for i, tab := range m.tabs {
		title := tab.Title()
		if len(title) > 20 {
			title = title[:17] + "..."
		}
		
		if tab.Modified() {
			title = modifiedIndicator + " " + title
		}
		
		var rendered string
		if tab.Active() {
			rendered = activeStyle.Render(title)
		} else {
			rendered = inactiveStyle.Render(title)
		}
		
		tabWidth := lipgloss.Width(rendered)
		if usedWidth+tabWidth > width-10 {
			// Show indicator that there are more tabs
			b.WriteString(inactiveStyle.Render("..."))
			break
		}
		
		b.WriteString(rendered)
		usedWidth += tabWidth
		
		if i < len(m.tabs)-1 {
			b.WriteString(" ")
			usedWidth++
		}
	}
	
	// Fill remaining space
	remaining := width - usedWidth
	if remaining > 0 {
		fillStyle := lipgloss.NewStyle().Background(lipgloss.Color("235"))
		b.WriteString(fillStyle.Render(strings.Repeat(" ", remaining)))
	}
	
	return b.String()
}
