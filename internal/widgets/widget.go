package widgets

// WIdget is the base interface for all widgets
type Widget interface {
	Show()
	Hide()
	IsVisible() bool
	Render() string
}
