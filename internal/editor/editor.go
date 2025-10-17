package editor

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/tobibamidele/minra/internal/buffer"
	"github.com/tobibamidele/minra/internal/clipboard"
	"github.com/tobibamidele/minra/internal/search"
	"github.com/tobibamidele/minra/internal/sidebar"
	"github.com/tobibamidele/minra/internal/statusbar"
	"github.com/tobibamidele/minra/internal/syntax"
	"github.com/tobibamidele/minra/internal/tabs"
	"github.com/tobibamidele/minra/internal/ui"
	"github.com/tobibamidele/minra/internal/viewport"
	"github.com/tobibamidele/minra/internal/widgets"
)

// Editor is the main editor model
type Editor struct {
	bufferMgr    *buffer.Manager
	tabMgr       *tabs.Manager
	clipboard    clipboard.Clipboard
	sidebar      *sidebar.Sidebar
	statusBar    *statusbar.StatusBar
	viewport     *viewport.Viewport
	highlighter  *syntax.Highlighter
	searchEngine *search.Engine
	renameWidget *widgets.RenameWidget
	searchWidget *widgets.SearchWidget
	mode         viewport.Mode
	width        int
	height       int
	statusMsg    string
	rootDir      string
}

// New creates a new editor
func New(rootDir string, config *Config) (*Editor, error) {
	bufferMgr := buffer.NewManager()
	tabMgr := tabs.NewManager()

	// Create initial buffer
	buf := bufferMgr.NewBuffer()
	tabMgr.NewTab(buf.ID(), "untitled")

	// Create sidebar
	sb, err := sidebar.New(rootDir, 35, 24)
	if err != nil {
		sb = nil
	}

	return &Editor{
		bufferMgr:    bufferMgr,
		tabMgr:       tabMgr,
		clipboard:    clipboard.New(),
		sidebar:      sb,
		statusBar:    statusbar.New(),
		viewport:     viewport.New(buf, viewport.ScreenWidth(), viewport.ScreenHeight()),
		highlighter:  syntax.New(),
		searchEngine: search.NewEngine(),
		renameWidget: widgets.NewRenameWidget(),
		searchWidget: widgets.NewSearchWidget(),
		mode:         viewport.ModeNormal,
		statusMsg:    "Press 'i' for insert mode, 'e' for sidebar, Ctrl+S to save",
		rootDir:      rootDir,
	}, nil
}

// Init initializes the editor
func (e *Editor) Init() tea.Cmd {
	lipgloss.SetColorProfile(termenv.TrueColor)
	return nil
}

// Update handles messages
func (e *Editor) Update(msg tea.Msg) (*Editor, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		e.width = msg.Width
		e.height = msg.Height

		// Update sidebar height
		if e.sidebar != nil {
			e.sidebar.SetHeight(e.height - 3)
		}

		// Update viewport size
		viewportWidth := e.getViewportWidth()
		viewportHeight := e.getViewportHeight()
		e.viewport.SetSize(viewportWidth, viewportHeight)

		return e, nil

	case tea.KeyMsg:
		return e, e.HandleKeyPress(msg)
	}

	return e, nil
}

// View renders the editor
func (e *Editor) View() string {
	// Render tabs
	tabBar := e.tabMgr.Render(e.width)

	// Render sidebar
	sidebarView := ""
	if e.sidebar != nil && e.sidebar.IsVisible() {
		sidebarView = e.sidebar.Render()
	}

	// Render viewport
	buf := e.bufferMgr.ActiveBuffer()
	var viewportView string
	if buf != nil {
		viewportView = e.viewport.Render(e.highlighter, buf.Cursor(), e.mode)
	}

	// // Wrap viewport in border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(e.getViewportWidth() + 5).
		Height(e.getViewportHeight()).
		Background(ui.ColorBackground)
	viewportView = borderStyle.Render(viewportView)

	// Combine sidebar and viewport
	mainView := ""
	if sidebarView != "" {
		mainView = lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, viewportView)
	} else {
		mainView = viewportView
	}

	// Overlay widgets if visible
	if e.renameWidget.IsVisible() {
		mainView = e.overlayWidget(mainView, e.renameWidget.Render(e.sidebar.Width()-5))
	}
	if e.searchWidget.IsVisible() {
		mainView = e.overlayWidget(mainView, e.searchWidget.Render())
	}

	// Render status bar
	statusBarView := e.renderStatusBar()

	// Combine everything
	return tabBar + "\n" + mainView + "\n" + statusBarView
}

func (e *Editor) getViewportWidth() int {
	sidebarWidth := 0
	if e.sidebar != nil && e.sidebar.IsVisible() {
		sidebarWidth = e.sidebar.Width()
	}
	return e.width - sidebarWidth - 6
}

func (e *Editor) getViewportHeight() int {
	return e.height - 4 // tabs + status bar + borders
}

func (e *Editor) renderStatusBar() string {
	buf := e.bufferMgr.ActiveBuffer()

	leftStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("240")).
		Padding(0, 1)

	modeStr := e.mode.String()
	left := leftStyle.Render(modeStr) + " " + e.statusMsg

	rightStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("240")).
		Padding(0, 1)

	modified := ""
	filename := "untitled"
	line := 1
	col := 1

	if buf != nil {
		if buf.Modified() {
			modified = "[+] "
		}
		if buf.Filepath() != "" {
			filename = filepath.Base(buf.Filepath())
		}
		cur := buf.Cursor()
		line = cur.Line() + 1
		col = cur.Col() + 1
	}

	right := rightStyle.Render(fmt.Sprintf("%s%s  Ln %d, Col %d",
		modified, filename, line, col))

	gap := e.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 0 {
		gap = 0
	}

	statusBar := lipgloss.NewStyle().
		Background(lipgloss.Color("240")).
		Width(e.width).
		Render(left + strings.Repeat(" ", gap) + right)

	return statusBar
}

func (e *Editor) overlayWidget(mainView, widgetView string) string {
	mainLines := strings.Split(mainView, "\n")
	widgetLines := strings.Split(widgetView, "\n")

	widgetWidth := 0
	for _, line := range widgetLines {
		if len(line) > widgetWidth {
			widgetWidth = len(line)
		}
	}

	widgetStartCol := e.width - widgetWidth - 2
	if widgetStartCol < 0 {
		widgetStartCol = 0
	}

	widgetStartRow := 1

	for i, widgetLine := range widgetLines {
		targetRow := widgetStartRow + i

		if targetRow >= len(mainLines) {
			mainLines = append(mainLines, strings.Repeat(" ", widgetStartCol)+widgetLine)
			continue
		}

		mainLine := mainLines[targetRow]
		mainRunes := []rune(mainLine)
		widgetRunes := []rune(widgetLine)

		var combined []rune

		if widgetStartCol < len(mainRunes) {
			combined = append(combined, mainRunes[:widgetStartCol]...)
		} else {
			combined = append(combined, mainRunes...)
			for len(combined) < widgetStartCol {
				combined = append(combined, ' ')
			}
		}

		combined = append(combined, widgetRunes...)

		widgetEndCol := widgetStartCol + len(widgetRunes)
		if widgetEndCol < len(mainRunes) {
			combined = append(combined, mainRunes[widgetEndCol:]...)
		}

		mainLines[targetRow] = string(combined)
	}

	return strings.Join(mainLines, "\n")
}
