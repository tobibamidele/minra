package buffer

// Change represents a buffer change for undo/redo
type Change struct {
	Line    int
	Col     int
	Text    string
	IsInsert bool
}

// History manages undo/redo
type History struct {
	changes []Change
	current int
	maxSize int
}

// NewHistory creates a new history
func NewHistory() *History {
	return &History{
		changes: make([]Change, 0),
		current: -1,
		maxSize: 100,
	}
}

// Record adds a change to history
func (h *History) Record(change Change) {
	// Remove any changes after current position
	if h.current < len(h.changes)-1 {
		h.changes = h.changes[:h.current+1]
	}
	
	h.changes = append(h.changes, change)
	h.current++
	
	// Limit history size
	if len(h.changes) > h.maxSize {
		h.changes = h.changes[1:]
		h.current--
	}
}

// Undo returns the change to undo
func (h *History) Undo() *Change {
	if h.current < 0 {
		return nil
	}
	
	change := h.changes[h.current]
	h.current--
	return &change
}

// Redo returns the change to redo
func (h *History) Redo() *Change {
	if h.current >= len(h.changes)-1 {
		return nil
	}
	
	h.current++
	return &h.changes[h.current]
}

// CanUndo returns if undo is possible
func (h *History) CanUndo() bool {
	return h.current >= 0
}

// CanRedo returns if redo is possible
func (h *History) CanRedo() bool {
	return h.current < len(h.changes)-1
}
