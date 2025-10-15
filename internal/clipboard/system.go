package clipboard

import (
	"github.com/atotto/clipboard"
)

// System uses system clipboard
type System struct {
	fallback *Internal
}

// NewSystem creates system clipboard
func NewSystem() *System {
	return &System{
		fallback: NewInternal(),
	}
}

func (c *System) Copy(text string) error {
	err := clipboard.WriteAll(text)
	if err != nil {
		return c.fallback.Copy(text)
	}
	c.fallback.Copy(text)
	return nil
}

func (c *System) Paste() (string, error) {
	text, err := clipboard.ReadAll()
	if err != nil || text == "" {
		return c.fallback.Paste()
	}
	return text, nil
}

func (c *System) Clear() error {
	clipboard.WriteAll("")
	return c.fallback.Clear()
}
