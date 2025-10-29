package filesearch

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tobibamidele/minra/internal/icons"
	"github.com/tobibamidele/minra/internal/ui"
)

type Sidebar struct {
	engine       *Engine
	results      []Result
	selected     int
	scrollOffset int
	visible      bool
	width, height int
	Query        string
	limit        int
}

func NewSidebar(engine *Engine, width, height int) *Sidebar {
	return &Sidebar{
		engine:       engine,
		results:      nil,
		selected:     0,
		scrollOffset: 0,
		visible:      true,
		width:        width,
		height:       height,
		Query:        "",
		limit:        500,
	}
}

func (s *Sidebar) SetSize(w, h int) { s.width = w; s.height = h }

func (s *Sidebar) Show() {
	s.visible = true
	s.selected = 0
	s.scrollOffset = 0
}

func (s *Sidebar) Hide() { s.visible = false }

func (s *Sidebar) IsVisible() bool { return s.visible }

func (s *Sidebar) UpdateQuery(q string) {
	s.Query = q
	if s.engine != nil && strings.TrimSpace(q) != "" {
		s.results = s.engine.Search(q, s.limit)
	} else {
		s.results = nil
	}
	s.selected = 0
	s.scrollOffset = 0
}

func (s *Sidebar) MoveUp() {
	if s.selected > 0 {
		s.selected--
		s.adjustScroll()
	}
}

func (s *Sidebar) MoveDown() {
	if s.selected < len(s.results)-1 {
		s.selected++
		s.adjustScroll()
	}
}

func (s *Sidebar) Selected() *Result {
	if len(s.results) == 0 {
		return nil
	}
	if s.selected < 0 || s.selected >= len(s.results) {
		return nil
	}
	return &s.results[s.selected]
}

func (s *Sidebar) adjustScroll() {
	visibleLines := s.height - 2
	if visibleLines <= 0 {
		return
	}
	if s.selected >= s.scrollOffset+visibleLines {
		s.scrollOffset = s.selected - visibleLines + 1
	}
	if s.selected < s.scrollOffset {
		s.scrollOffset = s.selected
	}
}

func (s *Sidebar) Render() string {
	if !s.visible {
		return ""
	}

	var b strings.Builder
	visibleLines := s.height - 2
	if visibleLines < 0 {
		visibleLines = 0
	}

	selectedBg := lipgloss.Color("240")
	selectedText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(selectedBg)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Background(ui.ColorSidebar).
		Bold(true).
		Width(s.width-2).
		Align(lipgloss.Center)

	baseStyle := lipgloss.NewStyle().Background(lipgloss.Color(ui.ColorSidebar))

	b.WriteString(titleStyle.Render("File Search") + "\n")

	endIdx := min(s.scrollOffset+visibleLines, len(s.results))
	for i := s.scrollOffset; i < endIdx; i++ {
		r := s.results[i]
		isSelected := i == s.selected


		// filename left, parent dir right-aligned
		fi := icons.GetFileIcon(r.Name)
		icon, iconColor := fi.Glyph, fi.Color
		iconStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(iconColor))

		name := " " + r.Name
		var nameStyled string
		if isSelected {
			nameStyled = selectedText.Render(name)
			iconStyle = iconStyle.Background(selectedBg)
		} else {
			nameStyled = lipgloss.NewStyle().Background(ui.ColorSidebar).Render(name)
		}

		nameLen := lipgloss.Width(name)
		if nameLen > 14 {
			name = " " + r.Name[:12] + ".."
			nameLen = lipgloss.Width(name)
		}
		
		parent := ""
		if r.ParentDir != "." && r.ParentDir != "" {
			parent = fmt.Sprintf("(%s)", r.ParentDir)
		}

		// compute padding so parent is right aligned
		// reserve 2 chars margin
		totalWidth := s.width - 4
		parentLen := lipgloss.Width(parent)
		spaceBetween := 1
		padding := 0
		if parentLen > 0 {
			padding = totalWidth - nameLen - parentLen - spaceBetween 
			if padding < 1 {
				padding = 1
			}
		} else {
			padding = 0
		}


		lineStyle := lipgloss.NewStyle().Width(s.width - 2)
		//
		// line := lipgloss.JoinHorizontal(
		// 	0,
		// 	iconStyle.Render(icon),
		// 	nameStyled,
		// 	lineStyle.Render(strings.Repeat(" ", padding)),
		// 	// lineStyle.Render(strings.Repeat(" ", spaceBetween)),
		// 	parent,
		// )
		var line strings.Builder
		line.WriteString(" " + iconStyle.Render(icon) + nameStyled)
		// Set the base style to match is isSelected
		if isSelected {
			baseStyle = baseStyle.Background(selectedBg)
		}

		line.WriteString(
			baseStyle.Render(strings.Repeat(" ", padding)) + 
			baseStyle.Render(strings.Repeat(" ", spaceBetween)) + baseStyle.Render(parent))

		if isSelected {
			temp := line
			line.Reset()
			line.WriteString(selectedText.Render(temp.String()))
			lineStyle = lineStyle.Background(selectedBg)
		} else {
			lineStyle = lineStyle.Background(ui.ColorSidebar)
		}

		b.WriteString(lineStyle.Render(line.String()) + "\n")
	}

	// pad rest
	for i := endIdx - s.scrollOffset; i < visibleLines; i++ {
		b.WriteString(strings.Repeat(" ", s.width-2) + "\n")
	}

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Background(ui.ColorSidebar).
		Width(s.width-2).
		Height(s.height-2)

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

