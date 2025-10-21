package sidebar

import (
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Icon struct {
	Glyph string
	Color string
}

var IconRegistry = struct {
	Extensions map[string]Icon
	Filenames  map[string]Icon
	Git        map[string]Icon
	Folders    map[string]Icon
}{
	Extensions: map[string]Icon{
		".go":			 {"\ue627", "39"},
		".py":			 {"\ue606", "226"},
		".js":			 {"\uf2ee", "220"},
		".ts":			 {"\ue69d", "33"},
		".tsx":			 {"\ued46", "33"},
		".jsx":			 {"\ued46", "220"},
		".json":		 {"\ue60b", "220"},
		".java":		 {"\ue256", "208"},
		".rs":			 {"\ue7a8", "208"},
		".rb":			 {"\ue21e", "197"},
		".r":			 {"\ue881", "39"},
		".swift":		 {"\ue755", "202"},
		".php":			 {"\ue608", "129"},
		".html":		 {"\ue736", "202"},
		".css":			 {"\ue749", "39"},
		".sh":			 {"\ue795", "34"},
		".yaml":		 {"\ue8eb", "33"},
		".yml":			 {"\ue8eb", "33"},
		".toml":		 {"\ue6b2", "250"},
		".ini":			 {"\ue615", "250"},
		".lock":		 {"\uf023", "244"},
		".log":			 {"\uf18d", "244"},
		".md":			 {"\ue73e", "244"},
		".pdf":			 {"\ue67d", "197"},
		".docx":		 {"\ue6a5", "33"},
		".docm":		 {"\ue6a5", "33"},
	},

	Filenames: map[string]Icon{
		".gitignore":     {"\ue702", "196"},
		".gitattributes": {"\ue702", "196"},
		"commit_editmsg": {"\ue702", "196"},
		".env":           {"\ue702", "250"},
		"makefile":       {"\ue673", "244"},
		"dockerfile":     {"\uf308", "33"},
		"license":        {"\ue60a", "250"},
		"go.mod":         {"\ue627", "39"},
		"go.sum":         {"\ue627", "39"},
		"cargo.toml":     {"\ue7a8", "208"},
		"package.json":   {"\ue60b", "220"},
		"requirements.txt": {"\ue606", "226"},
	},

	Git: map[string]Icon{
		"branch": {"\uf418", "39"}, // 
		"merge":  {"\ue727", "208"}, // 
		"tag":    {"\uf02b", "220"}, // 
		"stash":  {"\uf01c", "244"},
		"detached": {"\uf126", "244"},
	},

	Folders: map[string]Icon{
		"root":  {"\ue5fc", "250"},
		"git":   {"\ue5fb", "196"},
		"open":  {"\ue5fe", "250"},
		"closed":{"\ue5ff", "250"},
		"empty": {"\ue5ff", "250"},
	},
}

func GetFileIcon(filename string) Icon {
	ext := strings.ToLower(filepath.Ext(filename))
	name := strings.ToLower(filename)

	if icon, ok := IconRegistry.Extensions[ext]; ok {
		return icon
	}
	if icon, ok := IconRegistry.Filenames[name]; ok {
		return icon
	}
	return Icon{"\ue612", "250"} // default icon
}

func GetDefaultFileIcon() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	return style.Render("\ue612")
}

func GetDirectoryIcon(node *FileNode) string {
	downChevron := "\ueab4"
	leftChevron := "\ueab6"

	icon := downChevron + " " + IconRegistry.Folders["open"].Glyph

	if node.Name == ".git" && !node.Expanded {
		icon = leftChevron + " " + IconRegistry.Folders["git"].Glyph
	}
	if !node.Expanded {
		icon = leftChevron + " " + IconRegistry.Folders["closed"].Glyph
	}

	return icon
}
