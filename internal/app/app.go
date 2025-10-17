package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobibamidele/minra/internal/editor"
)

// App is the application coordinator
type App struct {
	editor *editor.Editor
	config *editor.Config
}

// New creates a new application
func New(rootDir string) (*App, error) {
	config := editor.DefaultConfig()

	ed, err := editor.New(rootDir, config)
	if err != nil {
		return nil, err
	}

	return &App{
		editor: ed,
		config: config,
	}, nil
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return a.editor.Init()
}

// Update handles messages
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.editor, cmd = a.editor.Update(msg)
	return a, cmd
}

// View renders the application
func (a *App) View() string {
	return a.editor.View()
}
