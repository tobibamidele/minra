package editor

// KeyType defines constants for keyboard input strings used across the editor.
// These map directly to the strings returned by tea.KeyMsg.String().
type KeyType string

const (
	// --- Global keys ---
	KeyQuit      KeyType = "ctrl+q"
	KeyInterrupt KeyType = "ctrl+c"
	KeySave      KeyType = "ctrl+s"
	KeyOpen      KeyType = "ctrl+o"
	KeyNew       KeyType = "ctrl+n"
	KeySidebar   KeyType = "ctrl+b"
	KeyNextBuf   KeyType = "alt+>"
	KeyPrevBuf   KeyType = "alt+<"
	KeyPaste     KeyType = "ctrl+i" // The terminal intercepts ctrl+v for paste so we don't get a proper key event.

	// --- Movement ---
	KeyUp       KeyType = "up"
	KeyDown     KeyType = "down"
	KeyLeft     KeyType = "left"
	KeyRight    KeyType = "right"
	KeyPageUp   KeyType = "pgup"
	KeyPageDown KeyType = "pgdown"
	KeyHome     KeyType = "home"
	KeyEnd      KeyType = "end"

	// --- Mode switches ---
	KeyInsert      KeyType = "i"
	KeySidebarMode KeyType = "e"
	KeyEscape      KeyType = "esc"

	// --- Editing ---
	KeyEnter     KeyType = "enter"
	KeyBackspace KeyType = "backspace"
	KeyTab       KeyType = "tab"
	KeyDelete    KeyType = "delete"

	// --- Navigation (vim-style) ---
	KeyH    KeyType = "h"
	KeyJ    KeyType = "j"
	KeyK    KeyType = "k"
	KeyL    KeyType = "l"
	KeyW    KeyType = "w"
	KeyB    KeyType = "b"
	KeyG    KeyType = "g"
	KeyBigG KeyType = "G"

	// --- Clipboard ---
	KeyY KeyType = "y"
	KeyP KeyType = "p"

	// --- Search / Rename ---
	KeySlash KeyType = "/"
	KeyR     KeyType = "r"

	// --- Regular keys ---
	Key0      KeyType = "0"
	KeyDollar KeyType = "$"
)

func (k KeyType) String() string {
	return string(k)
}
