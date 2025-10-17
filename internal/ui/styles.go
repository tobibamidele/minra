package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Border styles
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	// Selection styles
	SelectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230"))

	InactiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	// Text styles
	BoldStyle = lipgloss.NewStyle().
			Bold(true)

	ItalicStyle = lipgloss.NewStyle().
			Italic(true)

	// Status bar styles
	StatusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("240")).
			Foreground(lipgloss.Color("230"))
)
