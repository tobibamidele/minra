package tabs

// Tab represents a single tab
type Tab struct {
	id       string
	bufferID string
	title    string
	modified bool
	active   bool
}

// New creates a new tab
func New(bufferID, title string) *Tab {
	return &Tab{
		bufferID: bufferID,
		title:    title,
		modified: false,
		active:   false,
	}
}

// ID returns tab ID
func (t *Tab) ID() string {
	return t.id
}

// SetID sets tab ID
func (t *Tab) SetID(id string) {
	t.id = id
}

// BufferID returns associated buffer ID
func (t *Tab) BufferID() string {
	return t.bufferID
}

// Title returns tab title
func (t *Tab) Title() string {
	return t.title
}

// SetTitle sets tab title
func (t *Tab) SetTitle(title string) {
	t.title = title
}

// Modified returns if tab is modified
func (t *Tab) Modified() bool {
	return t.modified
}

// SetModified sets modified state
func (t *Tab) SetModified(modified bool) {
	t.modified = modified
}

// Active returns if tab is active
func (t *Tab) Active() bool {
	return t.active
}

// SetActive sets active state
func (t *Tab) SetActive(active bool) {
	t.active = active
}
