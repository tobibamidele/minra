package sidebar

import (
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var extToIconMap = map[string]string{
	".py":     "\ue606",
	".js":     "\ue74e",
	".ts":     "\ue628",
	".java":   "\ue256",
	".c":      "\ue61e",
	".h":      "\ue61e",
	".cpp":    "\ue61d",
	".hpp":    "\ue61e",
	".cs":     "\ue648",
	".go":     "\ue627",
	".rs":     "\ue7a8",
	".rb":     "\ue21e",
	".php":    "\ue608",
	".html":   "\ue736",
	".css":    "\ue749",
	".swift":  "\ue755",
	".kt":     "\ue634",
	".dart":   "\ue798",
	".scala":  "\ue737",
	".pl":     "\ue67a",
	".lua":    "\ue620",
	".sh":     "\ue795",
	".r":      "\ue67a",
	".elixir": "\ue62d",
	".md":     "\ue73e",
	".json":   "\ue60b",
}

var extToIconColor = map[string]string{
	".go":      "39",  // bright blue
	".py":      "226", // yellow
	".js":      "220", // golden yellow
	".json":    "220",
	".ts":      "33",  // light blue
	".java":    "208", // orange
	".rs":      "208", // rust orange
	".rb":      "197", // red
	".php":     "129", // purple
	".html":    "202", // orange
	".css":     "39",  // blue
	".dart":    "39",
	".swift":   "214", // orange-yellow
	".md":      "244", // gray
	".sh":      "34",  // green
	".c":       "45",  // cyan
	".cpp":     "45",
	"makefile": "244",
}

var filenameToIconMap = map[string]string{
	".gitignore":     "\ue702",
	".gitattributes": "\ue702",
	"commit_editmsg": "\ue702",
	"makefile":       "\ue673",
	"config":         "\ue702",
	".env":           "\ue702",
}

var filenameToColorMap = map[string]string{
	".gitignore":     "196",
	".gitattributes": "196",
	"commit_editmsg": "196",
	"makefile":       "244",
	"config":         "250",
	".env":           "250",
}

type dirIconsStruct struct {
	ClosedDirectory string
	OpenDirectory   string
	EmptyDirectory  string
	GitDirectory    string
	RootDirectory   string
}

var DirIcons = dirIconsStruct{
	ClosedDirectory: "\ue5ff",
	OpenDirectory:   "\ue5fe",
	EmptyDirectory:  "\ue5ff",
	GitDirectory:    "\ue5fb",
	RootDirectory:   "\ue5fc",
}

// GetFileIcon returns the PUA string of the file icon and the ANSI color code as a string
// Eg: test.go => ('\ue627', '39')
func GetFileIcon(filename string) (string, string) {
	// Get file extension. If not we continue with the filename as is
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = strings.ToLower(filename)
	}

	var icon string
	// We first check the extToIconMap
	icon, ok := extToIconMap[ext]
	if ok {
		color, ok := extToIconColor[ext]
		if ok {
			return icon, color // Return icon and color
		}
		return icon, "250" // Return the icon with the default color of grey
	}

	// Next we check the filenameToIconMap
	icon, ok = filenameToIconMap[ext]
	if ok {
		color, ok := filenameToColorMap[ext]
		if ok {
			return icon, color
		}
		return icon, "250"
	}

	// Return default file icon and color
	return "\ue612", "250"
}

func GetDefaultFileIcon() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	return style.Render("\ue612")
}

// GetDirectoryIcon returns the directory icon
func GetDirectoryIcon(node *FileNode) string {
	downChevron := "\ueab4"
	leftChevron := "\ueab6"

	icon := downChevron + " " + DirIcons.OpenDirectory

	if node.Name == ".git" && !node.Expanded {
		icon = leftChevron + " " + DirIcons.GitDirectory
	}

	if !node.Expanded {
		icon = leftChevron + " " + DirIcons.ClosedDirectory
	}

	return icon
}
