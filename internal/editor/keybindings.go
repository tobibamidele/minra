package editor

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobibamidele/minra/internal/viewport"
	"github.com/tobibamidele/minra/pkg/fileio"
)

// HandleKeyPress handles keyboard input
func (e *Editor) HandleKeyPress(msg tea.KeyMsg) tea.Cmd {
	// Global shortcuts
	switch msg.String() {
	case "ctrl+c", "ctrl+q":
		return tea.Quit
	case "ctrl+s":
		return e.SaveFile()
	case "ctrl+o":
		return e.openSelectedFile()
	case "ctrl+b":
		e.sidebar.Toggle()
		return nil
	case "ctrl+n":
		return e.NewFile()
	case "alt+>", "alt+.":
		e.NextBuffer()
		return nil
	case "alt+<", "alt+,":
		e.PreviousBuffer()
		return nil
	}

	// viewport.Mode specific handling
	switch e.mode {
	case viewport.ModeSidebar:
		return e.handleSidebarMode(msg)
	case viewport.ModeInsert:
		return e.handleInsertMode(msg)
	case viewport.ModeNormal:
		return e.handleNormalMode(msg)
	case viewport.ModeRename:
		return e.handleRenameMode(msg)
	case viewport.ModeSearch:
		return e.handleSearchMode(msg)
	}

	return nil
}

func (e *Editor) handleNormalMode(msg tea.KeyMsg) tea.Cmd {
	buf := e.bufferMgr.ActiveBuffer()
	if buf == nil {
		return nil
	}

	cur := buf.Cursor()

	switch msg.String() {
	case "i":
		e.mode = viewport.ModeInsert
		e.statusMsg = "-- INSERT --"
	case "e":
		if e.sidebar.IsVisible(){
			e.mode = viewport.ModeSidebar
			e.statusMsg = "-- SIDEBAR --"
		}
	case "h", "left":
		cur.MoveLeft(buf)
		e.viewport.AdjustScroll(cur)
	case "l", "right":
		cur.MoveRight(buf)
		e.viewport.AdjustScroll(cur)
		case "k", "up":
		cur.MoveUp(buf)
		e.viewport.AdjustScroll(cur)
	case "j", "down":
		cur.MoveDown(buf)
		e.viewport.AdjustScroll(cur)
	case "0", "home":
		cur.MoveToLineStart()
		e.viewport.AdjustScroll(cur)
	case "$", "end":
		cur.MoveToLineEnd(buf)
		e.viewport.AdjustScroll(cur)
	case "g":
		cur.MoveToBufferStart()
		e.viewport.AdjustScroll(cur)
	case "G":
		cur.MoveToBufferEnd(buf)
		e.viewport.AdjustScroll(cur)
	case "w":
		cur.MoveWordForward(buf)
		e.viewport.AdjustScroll(cur)
	case "b":
		cur.MoveWordBackward(buf)
		e.viewport.AdjustScroll(cur)
	case "y":
		// Copy line
		line := buf.Line(cur.Line())
		e.clipboard.Copy(line)
		e.statusMsg = "Copied line"
	case "p":
		// Paste
		text, _ := e.clipboard.Paste()
		if text != "" {
			e.pasteText(text)
			e.statusMsg = "Pasted"
		}
	case "/":
		e.mode = viewport.ModeSearch
		e.searchWidget.Show()
		e.statusMsg = "-- SEARCH --"
	}
	
	return nil
}

func (e *Editor) handleInsertMode(msg tea.KeyMsg) tea.Cmd {
	buf := e.bufferMgr.ActiveBuffer()
	if buf == nil {
		return nil
	}
	
	cur := buf.Cursor()
	
	switch msg.String() {
	case "esc":
		e.mode = viewport.ModeNormal
		e.statusMsg = "-- NORMAL --"
	case "backspace":
		buf.DeleteRune(cur.Line(), cur.Col())
		if cur.Col() > 0 {
			cur.MoveLeft(buf)
		} else if cur.Line() > 0 {
			cur.MoveUp(buf)
			cur.MoveToLineEnd(buf)
		}
		e.viewport.AdjustScroll(cur)
	case "enter":
		buf.InsertNewline(cur.Line(), cur.Col())
		cur.MoveDown(buf)
		cur.MoveToLineStart()
		e.viewport.AdjustScroll(cur)
	case "left":
		cur.MoveLeft(buf)
		e.viewport.AdjustScroll(cur)
	case "right":
		cur.MoveRight(buf)
		e.viewport.AdjustScroll(cur)
	case "up":
		cur.MoveUp(buf)
		e.viewport.AdjustScroll(cur)
	case "down":
		cur.MoveDown(buf)
		e.viewport.AdjustScroll(cur)
	case "home":
		cur.MoveToLineStart()
		e.viewport.AdjustScroll(cur)
	case "end":
		cur.MoveToLineEnd(buf)
		e.viewport.AdjustScroll(cur)
	case "tab":
		// Insert spaces for tab
		for i := 0; i < 4; i++ {
			buf.InsertRune(cur.Line(), cur.Col(), ' ')
			cur.MoveRight(buf)
		}
		e.viewport.AdjustScroll(cur)
	default:
		// Insert regular characters
		runes := []rune(msg.String())
		if len(runes) == 1 {
			buf.InsertRune(cur.Line(), cur.Col(), runes[0])
			cur.MoveRight(buf)
			e.viewport.AdjustScroll(cur)
		}
	}
	
	return nil
}

func (e *Editor) handleSidebarMode(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		e.mode = viewport.ModeNormal
		e.statusMsg = "-- NORMAL --"
	case "r":
		e.mode = viewport.ModeRename
		node := e.sidebar.SelectedNode()
		if node != nil {
			e.renameWidget.Show(node.Name)
			e.statusMsg = "-- RENAME --"
		}
	case "up", "k":
		e.sidebar.MoveUp()
	case "down", "j":
		e.sidebar.MoveDown()
	case "enter", " ":
		node := e.sidebar.SelectedNode()
		if node != nil {
			if node.IsDir {
				e.sidebar.ToggleSelected()
			} else {
				e.OpenFile(node.Path)
				e.mode = viewport.ModeNormal
			}
		}
	}
	
	return nil
}

func (e *Editor) handleRenameMode(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		e.renameWidget.Hide()
		e.mode = viewport.ModeSidebar
		e.statusMsg = "Cancelled"
	case "enter":
		newName := e.renameWidget.GetInput()
		if newName != "" {
			e.performRename(newName)
		}
		e.renameWidget.Hide()
		e.mode = viewport.ModeSidebar
	case "backspace":
		e.renameWidget.DeleteRune()
	case "left":
		e.renameWidget.MoveCursorLeft()
	case "right":
		e.renameWidget.MoveCursorRight()
	case "home":
		e.renameWidget.MoveCursorToStart()
	case "end":
		e.renameWidget.MoveCursorToEnd()
	default:
		runes := []rune(msg.String())
		if len(runes) == 1 {
			e.renameWidget.InsertRune(runes[0])
		}
	}
	
	return nil
}

func (e *Editor) handleSearchMode(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		e.searchWidget.Hide()
		e.mode = viewport.ModeNormal
		e.statusMsg = "-- NORMAL --"
	case "enter":
		query := e.searchWidget.GetInput()
		e.performSearch(query)
	case "backspace":
		e.searchWidget.DeleteRune()
	default:
		runes := []rune(msg.String())
		if len(runes) == 1 {
			e.searchWidget.InsertRune(runes[0])
		}
	}
	
	return nil
}

func (e *Editor) openSelectedFile() tea.Cmd {
	if e.sidebar == nil {
		return nil
	}

	node := e.sidebar.SelectedNode()
	if node != nil && !node.IsDir {
		return e.OpenFile(node.Path)
	}

	return nil
}

func (e *Editor) performRename(newName string) {
	node := e.sidebar.SelectedNode()
	if node == nil {
		return
	}

	oldPath := node.Path
	dir := filepath.Dir(oldPath)
	newPath := filepath.Join(dir, newName)

	if fileio.FileExists(newPath) && newPath != oldPath {
		e.statusMsg = fmt.Sprintf("File '%s' already exists", newName)
		return
	}

	err := fileio.RenameFile(oldPath, newPath)
	if err != nil {
		e.statusMsg = fmt.Sprintf("Error renaming: %v", err)
		return
	}

	e.sidebar.Refresh()
	e.statusMsg = fmt.Sprintf("Renamed to %s", newName)
}

func (e *Editor) performSearch(query string) {
	buf := e.bufferMgr.ActiveBuffer()
	if buf == nil {
		return
	}

	e.searchEngine.SetQuery(query)
	results := e.searchEngine.Search(buf)

	if len(results) > 0 {
		// Jump to first result
		cur := buf.Cursor()
		cur.SetPosition(results[0].Line, results[0].Column)
		e.viewport.AdjustScroll(cur)
		e.statusMsg = fmt.Sprintf("Found %d matches", len(results))
	} else {
		e.statusMsg = "No matches found"
	}

	e.searchWidget.Hide()
	e.mode = viewport.ModeNormal
}

func (e *Editor) pasteText(text string) {
	buf := e.bufferMgr.ActiveBuffer()
	if buf == nil {
		return
	}

	cur := buf.Cursor()
	buf.InsertText(cur.Line(), cur.Col(), text)
}
