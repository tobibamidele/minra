package editor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/tobibamidele/minra/internal/buffer"
	"github.com/tobibamidele/minra/internal/clipboard"
	fs "github.com/tobibamidele/minra/internal/filesearch"
	"github.com/tobibamidele/minra/internal/icons"
	"github.com/tobibamidele/minra/internal/search"
	"github.com/tobibamidele/minra/internal/sidebar"
	"github.com/tobibamidele/minra/internal/statusbar"
	"github.com/tobibamidele/minra/internal/syntax"
	"github.com/tobibamidele/minra/internal/tabs"
	"github.com/tobibamidele/minra/internal/ui"
	"github.com/tobibamidele/minra/internal/viewport"
	"github.com/tobibamidele/minra/internal/widgets"
	"github.com/tobibamidele/minra/pkg/utils"
)

// Editor is the main editor model
type Editor struct {
	bufferMgr				 *buffer.Manager
	tabMgr					 *tabs.Manager
	clipboard				 clipboard.Clipboard
	sidebar					 *sidebar.Sidebar
	statusBar				 *statusbar.StatusBar
	viewport				 *viewport.Viewport
	highlighter				 *syntax.Highlighter
	searchEngine			 *search.Engine
	renameWidget			 *widgets.RenameWidget
	searchWidget			 *widgets.SearchWidget
	fileSearchEngine		*fs.Engine
	fileSearchSidebar		*fs.Sidebar
	mode					 viewport.Mode
	width					 int
	height					 int
	statusMsg				 string
	rootDir					 string
}

// New creates a new editor
func New(rootDir string, config *Config) (*Editor, error) {
	bufferMgr := buffer.NewManager()
	tabMgr := tabs.NewManager()

	// Create initial buffer
	buf := bufferMgr.NewBuffer()
	tabMgr.NewTab(buf.ID(), "untitled")
	buf.SetTabSize(utils.GetTabSizeByFilePath(buf.Filepath()))

	// Create sidebar
	sbHeight, sbWidth := 35, 24
	sb, err := sidebar.New(rootDir, sbHeight, sbWidth)
	if err != nil {
		sb = nil
	}

	rootDir, _ = os.Getwd()
	engine, err := fs.NewEngine(rootDir)
	if err == nil {
		fmt.Printf("file count: %d", len(engine.Files()))
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
		fileSearchEngine: engine,
		fileSearchSidebar: fs.NewSidebar(engine, sb.Width(), viewport.ScreenHeight() - 3),
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

	// Render appropriate sidebar based on mode
	sidebarView := ""

	if e.mode == viewport.ModeFileSearch && e.fileSearchSidebar != nil && e.fileSearchSidebar.IsVisible() {
			// File search sidebar takes priority
			sidebarView = e.fileSearchSidebar.Render()
	} else if e.sidebar != nil && e.sidebar.IsVisible() {
			// Normal sidebar
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
		BorderBackground(ui.ColorBackground).
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
	leftChevron := "\ue0b0"  // Solid chevron (not \ue0b1)
	rightChevron := "\ue0b2" // Solid chevron (not \ue0b3)
	leftLineChevron := "\ue0b1"
	rightLineChevron := "\ue0b3"

	bgColor := lipgloss.Color("#252f3b")
	modeColor := lipgloss.Color("#9c9b9a")
	if e.mode == viewport.ModeInsert {
		modeColor = lipgloss.Color("#00ff00")
	}

	baseStyle := lipgloss.NewStyle().Background(bgColor)
	modeStyle := lipgloss.NewStyle().Background(modeColor)

	// Mode section
	modeStr := e.mode.String()

	// Chevron transition from mode to base (mode color fg, base color bg)
	modeChevronStyle := lipgloss.NewStyle().
		Foreground(modeColor).
		Background(bgColor)

	gitBranch, _ := utils.GetGitBranch(e.rootDir)
	gitBranchIcon, gitBranchIconColor := "", ""

	if gitBranch != "" {
		gitBranchIcon, gitBranchIconColor = icons.GetGitIcon("branch").Glyph, icons.GetGitIcon("branch").Color
	}

	gitBranchStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(gitBranchIconColor))
	// There's probably a better way to do this
	gitBranchStr := " " + gitBranchStyle.Render(gitBranchIcon) + baseStyle.Render(" ") + baseStyle.Render(gitBranchStyle.Render(gitBranch)) + baseStyle.Render(" ")
	noOfLinesAdded, err := e.numberOfLinesAdded()
	if err != nil {
		noOfLinesAdded = ""
	}

	left :=
	    modeStyle.Render(" "+lipgloss.NewStyle().Foreground(bgColor).Render(modeStr)) +
		modeChevronStyle.Render(leftChevron) +
		baseStyle.Render(gitBranchStr) +
		modeChevronStyle.Render(leftLineChevron)

	if noOfLinesAdded != "" {
		left = left + 
			baseStyle.Render(" " + noOfLinesAdded + " ") + modeChevronStyle.Render(leftLineChevron)
	}
		


	osIcon, _ := " "+ icons.GetOSIcon().Glyph+" ", icons.GetOSIcon().Color
	modified := ""
	filename := "untitled"
	fileType := ""
	line := 1
	col := 1
	if buf != nil {
		if buf.Modified() {
			modified = "[+] "
		}
		if buf.Filepath() != "" {
			filename = filepath.Base(buf.Filepath())
			fileType = " " + strings.Replace(filepath.Ext(filename), ".", "", 1) + " "
		}
		cur := buf.Cursor()
		line = cur.Line() + 1
		col = cur.Col() + 1
	}

	// Right side: filename < filetype < line:col
	rightText := fmt.Sprintf("%s%s ",
		modified, filename)

	lineColText := fmt.Sprintf(" %d:%d ", line, col)
	cursorLinePercent := " " + e.getCursorLinePercent() + " "
	fileEncoding, err := utils.DetectFileEncoding(e.bufferMgr.ActiveBuffer().Filepath())
	if err != nil {
		fileEncoding = "unknown"
	}

	fileEncoding = " " + fileEncoding + " "

	var right strings.Builder

	right.WriteString(
		baseStyle.Render(rightText) + 
		modeChevronStyle.Render(rightLineChevron) +
		baseStyle.Render(fileEncoding) +
		modeChevronStyle.Render(rightLineChevron) +
		baseStyle.Render(osIcon) +
		modeChevronStyle.Render(rightLineChevron))

	if fileType == "" {
		right.WriteString(
			baseStyle.Render(cursorLinePercent) +
			modeChevronStyle.Render(rightChevron) +
			modeStyle.Render(lipgloss.NewStyle().Foreground(bgColor).Render(lineColText)))
	} else {
		fileIconInfo := icons.GetFileIcon(filepath.Base(buf.Filepath()))
		fileIcon, fileIconColor := fileIconInfo.Glyph, fileIconInfo.Color
		fileIconStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(fileIconColor)).Background(bgColor)
		right.WriteString(
			baseStyle.Render(" " + fileIconStyle.Render(fileIcon) + baseStyle.Render(fileType)) +
			modeChevronStyle.Render(rightLineChevron) +
			baseStyle.Render(cursorLinePercent) +
			modeChevronStyle.Render(rightChevron) +
			modeStyle.Render(lipgloss.NewStyle().Foreground(bgColor).Render(lineColText)))
	}

	gap := e.width - lipgloss.Width(left) - lipgloss.Width(right.String())
	if gap < 0 {
		gap = 0
	}

	var statusBar strings.Builder
	statusBar.WriteString(lipgloss.NewStyle().
		Background(ui.ColorBackground).
		Width(e.width).
		Render(left + baseStyle.Render(strings.Repeat(" ", gap)) + right.String()))
	statusBar.WriteString("\n")
	statusBar.WriteString(lipgloss.NewStyle().Background(ui.ColorBackground).Width(e.width).Render(e.statusMsg))

	return statusBar.String()
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
