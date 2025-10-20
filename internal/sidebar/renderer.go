package sidebar

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tobibamidele/minra/internal/ui"
)

// Render renders the sidebar
func (s *Sidebar) Render() string {
	if !s.visible {
		return ""
	}

	var b strings.Builder
	flatList := s.tree.FlatList()
	visibleLines := s.height - 2 // Account for top and bottom borders

	// Styles
	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("240")).
		Foreground(lipgloss.Color("230"))

	dirStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true)

	// fileStyle := lipgloss.NewStyle().
	// 	Foreground(lipgloss.Color("252"))

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Width(s.width - 2).
		Align(lipgloss.Center)

	b.WriteString(titleStyle.Render("Files"))
	b.WriteString("\n")

	// Render visible lines
	endIdx := s.scrollOffset + visibleLines
	if endIdx > len(flatList) {
		endIdx = len(flatList)
	}

	for i := s.scrollOffset; i < endIdx; i++ {
		node := flatList[i]

		// Build the line with indentation
		indent := strings.Repeat("  ", node.Level)
		icon, iconColor := "", "250"
		if node.IsDir {
			icon = GetDirectoryIcon(node)
		} else {
			i := GetFileIcon(node.Name)
			icon, iconColor = i.Glyph, i.Color
		}

		lineText := indent + icon + " " + node.Name

		// Truncate if too long
		maxLen := s.width
		if len(lineText) > maxLen {
			lineText = lineText[:maxLen-2] + "..."
		}

		var style lipgloss.Style
		// Apply styling
		if i == s.selectedIndex {
			style = selectedStyle.Foreground(lipgloss.Color(iconColor))
		} else if node.IsDir {
			style = dirStyle
		} else {
			// line = fileStyle.Render(line)
			style = lipgloss.NewStyle().Foreground(lipgloss.Color(iconColor))
		}

		line := style.Width(s.width - 2).Render(lineText)
		b.WriteString(line)
		b.WriteString("\n")
	}

	// Fill remaining space with empty lines
	for i := endIdx - s.scrollOffset; i < visibleLines; i++ {
		b.WriteString(strings.Repeat(" ", s.width-2))
		b.WriteString("\n")
	}

	// Wrap in border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		// BorderBackground(ui.ColorBackground).
		Width(s.width - 2).
		Height(s.height - 2)
		// Background(ui.ColorBackground)

	return lipgloss.NewStyle().Background(ui.ColorBackground).Render(borderStyle.Render(b.String()))
}
