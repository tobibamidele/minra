package editor

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobibamidele/minra/internal/session"
	"github.com/tobibamidele/minra/internal/viewport"
	"github.com/tobibamidele/minra/pkg/fileio"
)


const BINARY_CONTENTS_HIDDEN_MSG string = "[binary file contents hidden]"
// SaveFile saves current buffer
func (e *Editor) SaveFile() tea.Cmd {
	buf := e.bufferMgr.ActiveBuffer()
	if buf == nil {
		e.statusMsg = "No buffer to save"
		return nil
	}

	if buf.Filepath() == "" {
		e.statusMsg = "No file path set"
		return nil
	}

	err := fileio.WriteFile(buf.Filepath(), buf.Content())
	if err != nil {
		e.statusMsg = fmt.Sprintf("Error saving: %v", err)
		return nil
	}

	e.mode = viewport.ModeNormal
	

	buf.SetModified(false)
	e.bufferMgr.ActiveBuffer().SetPreviousLineCount(e.bufferMgr.ActiveBuffer().LineCount())
	e.statusMsg = fmt.Sprintf("Saved: %s", filepath.Base(buf.Filepath()))
	return nil
}

// OpenFile opens a file
func (e *Editor) OpenFile(path string) tea.Cmd {
	content, err := fileio.ReadFile(path)
	if err != nil {
		e.statusMsg = fmt.Sprintf("Error opening: %v", err)
		return nil
	}

	isBinary := fileio.IsBinaryFile(path)

	buf, err := e.bufferMgr.OpenBuffer(path, content)
	if isBinary {
		buf, err = e.bufferMgr.OpenBinaryBuffer(path)
	}

	e.bufferMgr.ActiveBuffer().SetPreviousLineCount(e.bufferMgr.ActiveBuffer().LineCount())
	if err != nil {
		e.statusMsg = fmt.Sprintf("Error: %v", err)
		return nil
	}
	// Create tab for buffer
	e.tabMgr.NewTab(buf.ID(), filepath.Base(path))

	// Update viewport
	if isBinary {
		e.viewport.SetIsBinary(true)
	}
	e.viewport.SetBuffer(buf)

	// Detect language for syntax highlighting
	e.statusMsg = fmt.Sprintf("Opened: %s", filepath.Base(path))
	e.highlighter.ForExtension(filepath.Ext(path))
	return nil
}

// NewFile creates a new file
func (e *Editor) NewFile() tea.Cmd {
	buf := e.bufferMgr.NewBuffer()
	e.tabMgr.NewTab(buf.ID(), "untitled")
	e.viewport.SetBuffer(buf)
	e.statusMsg = "New file"
	return nil
}

// CloseFile closes the current file
func (e *Editor) CloseFile() tea.Cmd {
	buf := e.bufferMgr.ActiveBuffer()
	if buf == nil {
		return nil
	}

	if buf.Modified() {
		e.statusMsg = "File has unsaved changes"
		return nil
	}

	// Close the tab
	activeTab := e.tabMgr.ActiveTab()
	if activeTab != nil {
		e.tabMgr.CloseTab(activeTab.ID())
	}

	// Close the buffer
	e.bufferMgr.CloseBuffer(buf.ID())

	// Update viewport to new active buffer
	newBuf := e.bufferMgr.ActiveBuffer()
	if newBuf != nil {
		e.viewport.SetBuffer(newBuf)
	}

	e.statusMsg = "File closed"
	return nil
}

// NextBuffer swtiches to next buffer
func (e *Editor) NextBuffer() {
	e.bufferMgr.NextBuffer()
	e.tabMgr.NextTab()
	buf := e.bufferMgr.ActiveBuffer()
	if buf != nil {
		e.viewport.SetBuffer(buf)
	}
}

// PreviousBuffer switches to previous buffer
func (e *Editor) PreviousBuffer() {
	e.bufferMgr.PreviousBuffer()
	e.tabMgr.PreviousTab()
	buf := e.bufferMgr.ActiveBuffer()
	if buf != nil {
		e.viewport.SetBuffer(buf)
	}
}

// SaveState saves the current ui state
func (e *Editor) SaveState() error {
	return session.SaveUIState(e.sidebar, session.DefaultUIStatePath())
}

// LoadState loads the saved ui state
// func (e *Editor) LoadState() error {
//
// }
