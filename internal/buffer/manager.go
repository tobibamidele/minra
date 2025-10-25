package buffer

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tobibamidele/minra/pkg/utils"
)

// Manager manages multiple buffers
type Manager struct {
	buffers      map[string]*Buffer
	activeBuffer string
	bufferOrder  []string
}

// NewManager creates a new buffer manager
func NewManager() *Manager {
	return &Manager{
		buffers:     make(map[string]*Buffer),
		bufferOrder: make([]string, 0),
	}
}

// NewBuffer creates a new empty buffer
func (m *Manager) NewBuffer() *Buffer {
	id := uuid.New().String()
	buffer := New()
	buffer.SetID(id)

	m.buffers[id] = buffer
	m.bufferOrder = append(m.bufferOrder, id)
	m.activeBuffer = id

	return buffer
}

// OpenBuffer creates buffer from file
func (m *Manager) OpenBuffer(filepath, content string) (*Buffer, error) {
	// Check if already open
	for _, b := range m.buffers {
		if b.Filepath() == filepath {
			m.activeBuffer = b.ID()
			return b, nil
		}
	}

	id := uuid.New().String()
	buffer := NewFromContent(content, filepath)
	buffer.SetID(id)
	buffer.SetTabSize(utils.GetTabSizeByFilePath(filepath))

	m.buffers[id] = buffer
	m.bufferOrder = append(m.bufferOrder, id)
	m.activeBuffer = id

	return buffer, nil
}

// ActiveBuffer returns current buffer
func (m *Manager) ActiveBuffer() *Buffer {
	if m.activeBuffer == "" {
		return nil
	}
	return m.buffers[m.activeBuffer]
}

// CloseBuffer closes a buffer
func (m *Manager) CloseBuffer(id string) error {
	if _, ok := m.buffers[id]; !ok {
		return fmt.Errorf("buffer not found")
	}

	delete(m.buffers, id)

	for i, bufID := range m.bufferOrder {
		if bufID == id {
			m.bufferOrder = append(m.bufferOrder[:i], m.bufferOrder[i+1:]...)
			break
		}
	}

	if m.activeBuffer == id {
		if len(m.bufferOrder) > 0 {
			m.activeBuffer = m.bufferOrder[0]
		} else {
			m.activeBuffer = ""
		}
	}

	return nil
}

// NextBuffer switches to next buffer
func (m *Manager) NextBuffer() *Buffer {
	if len(m.bufferOrder) <= 1 {
		return m.ActiveBuffer()
	}

	idx := m.findBufferIndex(m.activeBuffer)
	nextIdx := (idx + 1) % len(m.bufferOrder)
	m.activeBuffer = m.bufferOrder[nextIdx]

	return m.ActiveBuffer()
}

// PreviousBuffer switches to previous buffer
func (m *Manager) PreviousBuffer() *Buffer {
	if len(m.bufferOrder) <= 1 {
		return m.ActiveBuffer()
	}

	idx := m.findBufferIndex(m.activeBuffer)
	prevIdx := (idx - 1 + len(m.bufferOrder)) % len(m.bufferOrder)
	m.activeBuffer = m.bufferOrder[prevIdx]

	return m.ActiveBuffer()
}

// AllBuffers returns all buffers in order
func (m *Manager) AllBuffers() []*Buffer {
	buffers := make([]*Buffer, 0, len(m.bufferOrder))
	for _, id := range m.bufferOrder {
		if buf, ok := m.buffers[id]; ok {
			buffers = append(buffers, buf)
		}
	}
	return buffers
}

// BufferCount returns number of buffers
func (m *Manager) BufferCount() int {
	return len(m.buffers)
}

func (m *Manager) findBufferIndex(id string) int {
	for i, bufID := range m.bufferOrder {
		if bufID == id {
			return i
		}
	}
	return -1
}

// MoveBufferLeft moves buffer left in order
func (m *Manager) MoveBufferLeft() {
	idx := m.findBufferIndex(m.activeBuffer)
	if idx <= 0 {
		return
	}
	m.bufferOrder[idx], m.bufferOrder[idx-1] = m.bufferOrder[idx-1], m.bufferOrder[idx]
}

// MoveBufferRight moves buffer right in order
func (m *Manager) MoveBufferRight() {
	idx := m.findBufferIndex(m.activeBuffer)
	if idx < 0 || idx >= len(m.bufferOrder)-1 {
		return
	}
	m.bufferOrder[idx], m.bufferOrder[idx+1] = m.bufferOrder[idx+1], m.bufferOrder[idx]
}
