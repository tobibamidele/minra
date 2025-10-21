package sidebar

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tobibamidele/minra/internal/ui"
)

// // Render renders the sidebar
// func (s *Sidebar) Render() string {
// 	if !s.visible {
// 		return ""
// 	}
//
// 	var b strings.Builder
// 	flatList := s.tree.FlatList()
// 	visibleLines := s.height - 2 // Account for top and bottom borders
//
// 	// Styles
// 	selectedStyle := lipgloss.NewStyle().
// 		Background(lipgloss.Color("240")).
// 		Foreground(lipgloss.Color("230"))
//
// 	dirStyle := lipgloss.NewStyle().
// 		Foreground(lipgloss.Color("39")).
// 		Bold(true)
//
// 	// fileStyle := lipgloss.NewStyle().
// 	// 	Foreground(lipgloss.Color("252"))
//
// 	// Title
// 	titleStyle := lipgloss.NewStyle().
// 		Foreground(lipgloss.Color("205")).
// 		Bold(true).
// 		Width(s.width - 2).
// 		Align(lipgloss.Center)
//
// 	b.WriteString(titleStyle.Render("Files"))
// 	b.WriteString("\n")
//
// 	// Render visible lines
// 	endIdx := s.scrollOffset + visibleLines
// 	if endIdx > len(flatList) {
// 		endIdx = len(flatList)
// 	}
//
// 	for i := s.scrollOffset; i < endIdx; i++ {
// 		node := flatList[i]
//
// 		// Build the line with indentation
// 		indent := strings.Repeat("  ", node.Level)
// 		icon, iconColor := "", "250"
// 		if node.IsDir {
// 			icon = GetDirectoryIcon(node)
// 		} else {
// 			i := GetFileIcon(node.Name)
// 			icon, iconColor = i.Glyph, i.Color
// 		}
//
// 		var iconStyle lipgloss.Style
// 		// Apply styling
// 		if i == s.selectedIndex {
// 			iconStyle = selectedStyle.Foreground(lipgloss.Color(iconColor))
// 		} else if node.IsDir {
// 			iconStyle = dirStyle
// 		} else {
// 			// line = fileStyle.Render(line)
// 			iconStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(iconColor))
// 		}
//
//
// 		iconText := indent + icon
//
// 		// Truncate if too long
// 		maxLen := s.width
// 		if len(iconText) > maxLen {
// 			iconText = iconText[:maxLen-2] + "..."
// 		}
//
// 		// line := style.Width(s.width - 2).Render(lineText)
// 		// Line we render the indent, icon with color and then the text
// 		// width := s.width - 3 - len(node.Name)
// 		line :=  iconStyle.Render(iconText) + " " + node.Name
// 		b.WriteString(line)
// 		b.WriteString("\n")
// 	}
//
// 	// Fill remaining space with empty lines
// 	for i := endIdx - s.scrollOffset; i < visibleLines; i++ {
// 		b.WriteString(strings.Repeat(" ", s.width-2))
// 		b.WriteString("\n")
// 	}
//
// 	// Wrap in border
// 	borderStyle := lipgloss.NewStyle().
// 		Border(lipgloss.NormalBorder()).
// 		BorderForeground(lipgloss.Color("240")).
// 		// BorderBackground(ui.ColorBackground).
// 		Width(s.width - 2).
// 		Height(s.height - 2)
// 		// Background(ui.ColorBackground)
//
// 	return lipgloss.NewStyle().Background(ui.ColorBackground).Render(borderStyle.Render(b.String()))
// }

// Render renders the sidebar
func (s *Sidebar) Render() string {
	if !s.visible {
		return ""
	}

	var b strings.Builder
	flatList := s.tree.FlatList()
	visibleLines := s.height - 2 // Account for top and bottom borders

	// Styles
	selectedLineStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("240")) // selected background

	selectedTextStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("240")).
		Foreground(lipgloss.Color("230"))// selected foreground (white)

	dirStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Width(s.width - 2).
		Align(lipgloss.Center)

	// Title
	b.WriteString(titleStyle.Render("Files") + "\n")

	endIdx := s.scrollOffset + visibleLines
	if endIdx > len(flatList) {
		endIdx = len(flatList)
	}

	// Render lines
	for i := s.scrollOffset; i < endIdx; i++ {
		node := flatList[i]

		indent := strings.Repeat("  ", node.Level)

		icon, iconColor := "", "250"
		if node.IsDir {
			icon = GetDirectoryIcon(node)
		} else {
			fi := GetFileIcon(node.Name)
			icon, iconColor = fi.Glyph, fi.Color
		}

		// icon always has its own color
		iconStyled := lipgloss.NewStyle().
			Foreground(lipgloss.Color(iconColor)).
			Render(icon)

		// text styling
		textStyled := node.Name
		isSelected := i == s.selectedIndex

		if isSelected {
			textStyled = selectedTextStyle.Render(node.Name)
		} else if node.IsDir {
			textStyled = dirStyle.Render(node.Name)
		}

		var fullLine string
		raw := indent + icon + " " + node.Name

		if isSelected {
			// Entire line gets background first
			fullLine = selectedLineStyle.Width(s.width - 2).Render(raw)

			// Re-apply foreground colors *inside* the styled line
			// fullLine = strings.Replace(fullLine, " ", selectedTextStyle.Render(" "), 1)
			fullLine = strings.Replace(fullLine, icon, iconStyled, 1)
			// Had a problem where it render the background and then skipped the space before the file name
			// Replace " " + fileName with the rendered version
			fullLine = strings.Replace(fullLine, fmt.Sprintf(" %s", node.Name), selectedTextStyle.Render(" " + node.Name), 1)
		} else {
			// normal case
			fullLine = indent + iconStyled + " " + textStyled
		}


		// fullLine := indent + iconStyled + " " + textStyled

		// If selected: apply background to full width
		if isSelected {
			fullLine = selectedLineStyle.Background(lipgloss.Color("240")).Width(s.width - 2).Render(fullLine)
		}

		b.WriteString(fullLine + "\n")
	}

	// Fill empty lines
	for i := endIdx - s.scrollOffset; i < visibleLines; i++ {
		b.WriteString(strings.Repeat(" ", s.width-2) + "\n")
	}

	// Wrap in border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(s.width - 2).
		Height(s.height - 2)

	return lipgloss.NewStyle().
		Background(ui.ColorBackground).
		Render(borderStyle.Render(b.String()))
}
