package clipboard

// Internal is in-memory clipboard
type Internal struct {
	content    string
	history    []string
	maxHistory int
}

// NewInternal creates internal clipboard
func NewInternal() *Internal {
	return &Internal{
		history:    make([]string, 0),
		maxHistory: 10,
	}
}

func (c *Internal) Copy(text string) error {
	c.content = text
	c.history = append(c.history, text)
	if len(c.history) > c.maxHistory {
		c.history = c.history[1:]
	}
	return nil
}

func (c *Internal) Paste() (string, error) {
	return c.content, nil
}

func (c *Internal) Clear() error {
	c.content = ""
	return nil
}

func (c *Internal) History() []string {
	return c.history
}
