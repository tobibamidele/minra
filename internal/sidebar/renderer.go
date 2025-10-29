package sidebar

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tobibamidele/minra/internal/icons"
	"github.com/tobibamidele/minra/internal/ui"
)

// Render renders the sidebar
func (s *Sidebar) Render() string {
	if !s.visible {
		return ""
	}

	var b strings.Builder
	flatList := s.tree.FlatList()
	visibleLines := s.height - 2 // minus top & bottom borders

	// Styles
	selectedBg := lipgloss.Color("240")
	selectedText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(selectedBg)

	dirStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Background(ui.ColorSidebar).
		Bold(true)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Background(ui.ColorSidebar).
		Bold(true).
		Width(s.width - 2).
		Align(lipgloss.Center)

	// Title
	b.WriteString(titleStyle.Render("Files") + "\n")

	endIdx := min(s.scrollOffset+visibleLines, len(flatList))

	for i := s.scrollOffset; i < endIdx; i++ {
		node := flatList[i]
		isSelected := i == s.selectedIndex

		indent := strings.Repeat("  ", node.Level)

		// icon
		icon, iconColor := "", "250"
		if node.IsDir {
			icon = icons.GetDirectoryIcon(node.Name, node.Expanded)
		} else {
			fi := icons.GetFileIcon(node.Name)
			icon, iconColor = fi.Glyph, fi.Color
		}
		iconStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(iconColor))

		// text styling
		var nameStyled string
		if isSelected {
			nameStyled = selectedText.Render(" " + node.Name)
			iconStyle = iconStyle.Background(selectedBg)
		} else if node.IsDir {
			nameStyled = dirStyle.Render(" " + node.Name)
		} else {
			nameStyled = lipgloss.NewStyle().Background(ui.ColorSidebar).Render(" " + node.Name)
		}

		// join parts (indent + icon + filename)
		line := lipgloss.JoinHorizontal(
			0,
			indent,
			iconStyle.Render(icon),
			nameStyled,
		)

		// final width + background
		lineStyle := lipgloss.NewStyle().
			Width(s.width - 2)

		if isSelected {
			lineStyle = lineStyle.Background(selectedBg)
		} else {
			lineStyle = lineStyle.Background(ui.ColorSidebar)
		}

		b.WriteString(lineStyle.Render(line) + "\n")
	}

	// pad the rest
	for i := endIdx - s.scrollOffset; i < visibleLines; i++ {
		b.WriteString(strings.Repeat(" ", s.width-2) + "\n")
	}

	// border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Background(ui.ColorSidebar).
		Width(s.width - 2).
		Height(s.height - 2)

	return lipgloss.NewStyle().
		Background(ui.ColorBackground).
		Render(borderStyle.Render(b.String()))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
