package sidebar

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileNode represents a file or directory
type FileNode struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"isDir"`
	Children []*FileNode `json:"children,omitempty"`
	Expanded bool        `json:"expanded"`
	Level    int         `json:"-"`
}

// FileTree represents the file tree
type FileTree struct {
	Root     *FileNode
	flatList []*FileNode
}

// NewFileTree creates a new file tree
func NewFileTree(rootPath string) (*FileTree, error) {
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	root := &FileNode{
		Name:     filepath.Base(absPath),
		Path:     absPath,
		IsDir:    true,
		Expanded: true,
		Level:    0,
	}

	if err := loadDirectory(root); err != nil {
		return nil, err
	}

	tree := &FileTree{Root: root}
	tree.rebuildFlatList()
	return tree, nil
}

func loadDirectory(node *FileNode) error {
	if !node.IsDir {
		return nil
	}

	entries, err := os.ReadDir(node.Path)
	if err != nil {
		return err
	}

	node.Children = make([]*FileNode, 0)

	for _, entry := range entries {
		// Skip hidden directories except .git
		if strings.HasPrefix(entry.Name(), ".") && !strings.HasPrefix(entry.Name(), ".git") {
			continue
		}

		child := &FileNode{
			Name:     entry.Name(),
			Path:     filepath.Join(node.Path, entry.Name()),
			IsDir:    entry.IsDir(),
			Expanded: false,
			Level:    node.Level + 1,
		}

		node.Children = append(node.Children, child)
	}

	sort.Slice(node.Children, func(i, j int) bool {
		if node.Children[i].IsDir != node.Children[j].IsDir {
			return node.Children[i].IsDir
		}
		return strings.ToLower(node.Children[i].Name) < strings.ToLower(node.Children[j].Name)
	})

	return nil
}

func (t *FileTree) rebuildFlatList() {
	t.flatList = make([]*FileNode, 0)
	t.addToFlatList(t.Root)
}

func (t *FileTree) addToFlatList(node *FileNode) {
	t.flatList = append(t.flatList, node)

	if node.IsDir && node.Expanded {
		for _, child := range node.Children {
			t.addToFlatList(child)
		}
	}
}

// FlatList returns visible nodes
func (t *FileTree) FlatList() []*FileNode {
	return t.flatList
}

// ToggleExpanded toggles node expansion
func (t *FileTree) ToggleExpanded(node *FileNode) error {
	if !node.IsDir {
		return nil
	}

	node.Expanded = !node.Expanded

	if node.Expanded && len(node.Children) == 0 {
		if err := loadDirectory(node); err != nil {
			node.Expanded = false
			return err
		}
	}

	t.rebuildFlatList()
	return nil
}

// Refresh reloads the tree
func (t *FileTree) Refresh() error {
	t.Root.Children = nil
	if err := loadDirectory(t.Root); err != nil {
		return err
	}
	t.rebuildFlatList()
	return nil
}
