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
	switch KeyType(msg.String()) {
	case KeyQuit, KeyInterrupt:
		switch e.mode {
		case viewport.ModeRename:
			e.renameWidget.Hide()
			e.mode = viewport.ModeSidebar
			e.statusMsg = "Cancelled"
			return nil
		case viewport.ModeSearch:
			e.searchWidget.Hide()
			e.mode = viewport.ModeSidebar
			e.statusMsg = "Cancelled"
			return nil
		default:
			// Attempt to save state
			e.SaveState()
			return tea.Quit
		}
	case KeySave:
		return e.SaveFile()
	case KeyOpen:
		return e.openSelectedFile()
	case KeySidebar:
		e.sidebar.Toggle()
		return nil
	case KeyNew:
		return e.NewFile()
	case KeyNextBuf, "alt+.":
		e.NextBuffer()
		return nil
	case KeyPrevBuf, "alt+,":
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

	switch KeyType(msg.String()) {
	case KeyInsert:
		e.mode = viewport.ModeInsert
		e.statusMsg = "-- INSERT --"
	case KeySidebarMode:
		if e.sidebar.IsVisible() {
			e.mode = viewport.ModeSidebar
			e.statusMsg = "-- SIDEBAR --"
		}
	case KeyH, KeyLeft:
		cur.MoveLeft(buf)
		e.viewport.AdjustScroll(cur)
	case KeyL, KeyRight:
		cur.MoveRight(buf)
		e.viewport.AdjustScroll(cur)
	case KeyK, KeyUp:
		cur.MoveUp(buf)
		e.viewport.AdjustScroll(cur)
	case KeyJ, KeyDown:
		cur.MoveDown(buf)
		e.viewport.AdjustScroll(cur)
	case Key0, KeyHome:
		cur.MoveToLineStart()
		e.viewport.AdjustScroll(cur)
	case KeyDollar, KeyEnd:
		cur.MoveToLineEnd(buf)
		e.viewport.AdjustScroll(cur)
	case KeyBackspace:
		cur.MoveLeft(buf)
		e.viewport.AdjustScroll(cur)
	case KeyG:
		cur.MoveToBufferStart()
		e.viewport.AdjustScroll(cur)
	case KeyBigG:
		cur.MoveToBufferEnd(buf)
		e.viewport.AdjustScroll(cur)
	case KeyW:
		cur.MoveWordForward(buf)
		e.viewport.AdjustScroll(cur)
	case KeyB:
		cur.MoveWordBackward(buf)
		e.viewport.AdjustScroll(cur)
	case KeyPageDown:
		cur.MovePageDown(buf, e.getViewportHeight())
		e.viewport.AdjustScroll(cur)
	case KeyPageUp:
		cur.MovePageUp(buf, e.getViewportHeight())
		e.viewport.AdjustScroll(cur)
	case KeyY:
		// Copy line
		line := buf.Line(cur.Line())
		e.clipboard.Copy(line)
		e.statusMsg = "Copied line"
	case KeyP:
		// Paste
		text, _ := e.clipboard.Paste()
		if text != "" {
			e.pasteText(text)
			e.statusMsg = "Pasted"
		}
	case KeySlash:
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

	switch KeyType(msg.String()) {
	case KeyEscape:
		e.mode = viewport.ModeNormal
		e.statusMsg = "-- NORMAL --"
	case KeyBackspace:
		buf.DeleteRune(cur.Line(), cur.Col())
		if cur.Col() > 0 {
			cur.MoveLeft(buf)
		} else if cur.Line() > 0 {
			cur.MoveUp(buf)
			cur.MoveToLineEnd(buf)
		}
		e.viewport.AdjustScroll(cur)
	case KeyEnter:
		buf.InsertNewline(cur.Line(), cur.Col())
		cur.MoveDown(buf)
		// cur.MoveToLineStart() // <== removed coz it sets the position in InsertNewline
		e.viewport.AdjustScroll(cur)
	case KeyLeft:
		cur.MoveLeft(buf)
		e.viewport.AdjustScroll(cur)
	case KeyRight:
		cur.MoveRight(buf)
		e.viewport.AdjustScroll(cur)
	case KeyUp:
		cur.MoveUp(buf)
		e.viewport.AdjustScroll(cur)
	case KeyDown:
		cur.MoveDown(buf)
		e.viewport.AdjustScroll(cur)
	case KeyHome:
		cur.MoveToLineStart()
		e.viewport.AdjustScroll(cur)
	case KeyEnd:
		cur.MoveToLineEnd(buf)
		e.viewport.AdjustScroll(cur)
	case KeyPageDown:
		cur.MovePageDown(buf, e.getViewportHeight())
		e.viewport.AdjustScroll(cur)
	case KeyPageUp:
		cur.MovePageUp(buf, e.getViewportHeight())
		e.viewport.AdjustScroll(cur)
	case KeyDelete:
		buf.DeleteLine(cur.Line()) // TODO: Move th cursor to the previous line
	case KeyPaste:
		// Paste from clipboard
		text, _ := e.clipboard.Paste()
		if text != "" {
			e.pasteText(text)
			e.statusMsg = "Pasted"
		}
	case "tab":
		// Insert spaces for tab
		for i := 0; i < e.viewport.TabSize(); i++ {
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
	case "ctrl+c":
		e.renameWidget.Hide()
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
