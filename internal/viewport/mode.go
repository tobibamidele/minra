package viewport

type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeVisual
	ModeSidebar
	ModeCommand
	ModeRename
	ModeSearch
)

func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "NORMAL"
	case ModeInsert:
		return "INSERT"
	case ModeSidebar:
		return "SIDEBAR"
	case ModeVisual:
		return "VISUAL"
	case ModeCommand:
		return "COMMAND"
	case ModeRename:
		return "RENAME"
	case ModeSearch:
		return "SEARCH"
	default:
		return "UNKNOWN"
	}
}
