# Minra - TUI Text Editor
## Description
Minra is a terminal based text editor written in go. It supports basic text editing capabilities like create files, find and replace, rename files among other.

## Getting started
*Ensure you have the go toolkit installed version 1.24+*

- Clone the repo
```bash
git clone https://github.com/tobibamidele/minra.git
```

- Compile and run
```bash
make build
./bin/minra
```

## File structure
```bash
minra/
├── main.go
├── Makefile
├── go.mod
├── go.sum
├── README.md
├── .gitignore
│
├── cmd/
│   └── minra/
│       └── main.go                    # Entry point (moved from root)
│
├── internal/                          # Private application code
│   ├── app/
│   │   ├── app.go                     # Main application coordinator
│   │   └── config.go                  # Application configuration
│   │
│   ├── editor/
│   │   ├── editor.go                  # Main editor model (simplified)
│   │   ├── modes.go                   # Mode definitions and handlers
│   │   ├── commands.go                # Command execution (save, open, etc.)
│   │   └── keybindings.go            # Centralized keybinding management
│   │
│   ├── buffer/
│   │   ├── buffer.go                  # Core buffer implementation
│   │   ├── manager.go                 # Buffer manager (multiple buffers)
│   │   ├── operations.go              # Buffer operations (insert, delete, etc.)
│   │   └── history.go                 # Undo/redo history per buffer
│   │
│   ├── cursor/
│   │   ├── cursor.go                  # Cursor implementation
│   │   └── movement.go                # Cursor movement logic
│   │
│   ├── clipboard/
│   │   ├── clipboard.go               # Clipboard interface
│   │   ├── system.go                  # System clipboard integration
│   │   ├── internal.go                # Internal clipboard (fallback)
│   │   └── operations.go              # Copy/cut/paste operations
│   │
│   ├── tabs/
│   │   ├── tab.go                     # Single tab representation
│   │   ├── manager.go                 # Tab manager (switching, ordering)
│   │   └── renderer.go                # Tab bar rendering
│   │
│   ├── viewport/
│   │   ├── viewport.go                # Viewport for rendering visible area
│   │   ├── scroll.go                  # Scrolling logic
│   │   └── renderer.go                # Main editor content renderer
│   │
│   ├── sidebar/
│   │   ├── sidebar.go                 # Sidebar model
│   │   ├── filetree.go               # File tree structure
│   │   ├── icons.go                   # File icons
│   │   └── renderer.go                # Sidebar rendering
│   │
│   ├── statusbar/
│   │   ├── statusbar.go               # Status bar model
│   │   └── renderer.go                # Status bar rendering
│   │
│   ├── widgets/
│   │   ├── widget.go                  # Base widget interface
│   │   ├── rename.go                  # Rename widget
│   │   ├── create.go                  # Create file widget
│   │   ├── search.go                  # Search/find widget
│   │   ├── command_palette.go         # Command palette widget
│   │   └── dialog.go                  # Generic dialog widget
│   │
│   ├── syntax/
│   │   ├── highlighter.go             # Syntax highlighter
│   │   ├── languages/
│   │   │   ├── go.go                  # Go language rules
│   │   │   ├── python.go              # Python language rules
│   │   │   ├── javascript.go          # JavaScript language rules
│   │   │   └── base.go                # Base language interface
│   │   └── theme.go                   # Color theme management
│   │
│   ├── search/
│   │   ├── search.go                  # Search engine
│   │   ├── replace.go                 # Find and replace
│   │   └── regex.go                   # Regex search support
│   │
│   ├── session/
│   │   ├── session.go                 # Session management
│   │   ├── persistence.go             # Save/restore session
│   │   └── workspace.go               # Workspace settings
│   │
│   └── ui/
│       ├── styles.go                  # Centralized UI styles
│       ├── colors.go                  # Color definitions
│       ├── layout.go                  # Layout calculations
│       └── components/
│           ├── border.go              # Border components
│           ├── scrollbar.go           # Scrollbar component
│           └── modal.go               # Modal overlay component
│
├── pkg/                               # Public, reusable packages
│   ├── fileio/
│   │   ├── reader.go                  # File reading utilities
│   │   ├── writer.go                  # File writing utilities
│   │   └── watcher.go                 # File system watcher
│   │
│   └── utils/
│       ├── strings.go                 # String utilities
│       ├── paths.go                   # Path utilities
│       └── encoding.go                # Encoding detection
│
├── configs/
│   ├── default.yaml                   # Default configuration
│   └── keybindings.yaml              # Default keybindings
│
└── docs/
    ├── architecture.md                # Architecture documentation
    ├── contributing.md                # Contribution guidelines
    └── features.md                    # Feature documentation
```

## TODO
- [X] Fix cursor
- [X] Fix tab identification - Currently doesn't identify tabs
- [ ] Fix ANSI escape codes messing up editor
- [ ] Add auto-completion suggestions
- [ ] Add file parser to allow collapsing blocks of a file
- [ ] Add create file support
- [X] Add find and replace support
- [ ] Add LSP support (way in the future)
