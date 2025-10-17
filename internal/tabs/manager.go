package tabs

import "github.com/google/uuid"

// Manager manages tabs
type Manager struct {
	tabs      []*Tab
	activeIdx int
	maxTabs   int
}

// NewManager creates new tab manager
func NewManager() *Manager {
	return &Manager{
		tabs:      make([]*Tab, 0),
		activeIdx: -1,
		maxTabs:   20,
	}
}

// NewTab creates a new tab
func (m *Manager) NewTab(bufferID, title string) *Tab {
	tab := New(bufferID, title)
	tab.SetID(uuid.New().String())

	m.tabs = append(m.tabs, tab)
	m.activeIdx = len(m.tabs) - 1
	m.updateActiveStates()

	return tab
}

// CloseTab closes a tab
func (m *Manager) CloseTab(id string) {
	for i, tab := range m.tabs {
		if tab.ID() == id {
			m.tabs = append(m.tabs[:i], m.tabs[i+1:]...)
			if m.activeIdx >= len(m.tabs) && len(m.tabs) > 0 {
				m.activeIdx = len(m.tabs) - 1
			}
			if len(m.tabs) == 0 {
				m.activeIdx = -1
			}
			m.updateActiveStates()
			return
		}
	}
}

// NextTab switches to next tab
func (m *Manager) NextTab() {
	if len(m.tabs) <= 1 {
		return
	}
	m.activeIdx = (m.activeIdx + 1) % len(m.tabs)
	m.updateActiveStates()
}

// PreviousTab switches to previous tab
func (m *Manager) PreviousTab() {
	if len(m.tabs) <= 1 {
		return
	}
	m.activeIdx = (m.activeIdx - 1 + len(m.tabs)) % len(m.tabs)
	m.updateActiveStates()
}

// ActiveTab returns active tab
func (m *Manager) ActiveTab() *Tab {
	if m.activeIdx < 0 || m.activeIdx >= len(m.tabs) {
		return nil
	}
	return m.tabs[m.activeIdx]
}

// AllTabs returns all tabs
func (m *Manager) AllTabs() []*Tab {
	return m.tabs
}

// TabCount returns number of tabs
func (m *Manager) TabCount() int {
	return len(m.tabs)
}

// MoveTabLeft moves active tab left
func (m *Manager) MoveTabLeft() {
	if m.activeIdx <= 0 {
		return
	}
	m.tabs[m.activeIdx], m.tabs[m.activeIdx-1] = m.tabs[m.activeIdx-1], m.tabs[m.activeIdx]
	m.activeIdx--
}

// MoveTabRight moves active tab right
func (m *Manager) MoveTabRight() {
	if m.activeIdx < 0 || m.activeIdx >= len(m.tabs)-1 {
		return
	}
	m.tabs[m.activeIdx], m.tabs[m.activeIdx+1] = m.tabs[m.activeIdx+1], m.tabs[m.activeIdx]
	m.activeIdx++
}

func (m *Manager) updateActiveStates() {
	for i, tab := range m.tabs {
		tab.SetActive(i == m.activeIdx)
	}
}
