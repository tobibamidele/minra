package statusbar

// StatusBar represents th status bar
type StatusBar struct {
	message string
}

// New creates a new status bar
func New() *StatusBar {
	return &StatusBar{
		message: "",
	}
}

// SetMessage sets status message
func (s *StatusBar) SetMessage(msg string) {
	s.message = msg
}

// Message returns current message
func (s *StatusBar) Message() string {
	return s.message
}
