package syntax

import "github.com/charmbracelet/lipgloss"

// Theme defines color scheme
type Theme struct {
	Keyword  lipgloss.Style
	Type     lipgloss.Style
	Constant lipgloss.Style
	String   lipgloss.Style
	Comment  lipgloss.Style
	Function lipgloss.Style
	Number   lipgloss.Style
}

// DefaultTheme returns default theme
func DefaultTheme() *Theme {
	return &Theme{
		Keyword:  lipgloss.NewStyle().Foreground(lipgloss.Color("205")),
		Type:     lipgloss.NewStyle().Foreground(lipgloss.Color("117")),
		Constant: lipgloss.NewStyle().Foreground(lipgloss.Color("141")),
		String:   lipgloss.NewStyle().Foreground(lipgloss.Color("150")),
		Comment:  lipgloss.NewStyle().Foreground(lipgloss.Color("243")),
		Function: lipgloss.NewStyle().Foreground(lipgloss.Color("220")),
		Number:   lipgloss.NewStyle().Foreground(lipgloss.Color("174")),
	}
}
