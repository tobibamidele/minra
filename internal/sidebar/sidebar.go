package sidebar

// Sidebar represents the file browser sidebar
type Sidebar struct {
	tree          *FileTree
	selectedIndex int
	visible       bool
	width         int
	height        int
	scrollOffset  int
}

// NewSidebar creates a new sidebar
func New(rootPath string, width, height int) (*Sidebar, error) {
	tree, err := NewFileTree(rootPath)
	if err != nil {
		return nil, err
	}

	return &Sidebar{
		tree:          tree,
		selectedIndex: 0,
		visible:       true,
		width:         width,
		height:        height,
		scrollOffset:  0,
	}, nil
}

// IsVisible returns whether the sidebar is visible
func (s *Sidebar) IsVisible() bool {
	return s.visible
}

// Toggle toggles the sidebar visibility
func (s *Sidebar) Toggle() {
	s.visible = !s.visible
}

// Width returns the sidebar width
func (s *Sidebar) Width() int {
	if !s.visible {
		return 0
	}
	return s.width
}

// SetHeight sets the sidebar height
func (s *Sidebar) SetHeight(height int) {
	s.height = height
}

// MoveUp moves the selection up
func (s *Sidebar) MoveUp() {
	if s.selectedIndex > 0 {
		s.selectedIndex--
		s.adjustScroll()
	}
}

// MoveDown moves the selection down
func (s *Sidebar) MoveDown() {
	flatList := s.tree.FlatList()
	if s.selectedIndex < len(flatList)-1 {
		s.selectedIndex++
		s.adjustScroll()
	}
}

// ToggleSelected toggles the expanded state of the selected node
func (s *Sidebar) ToggleSelected() error {
	flatList := s.tree.FlatList()
	if s.selectedIndex >= 0 && s.selectedIndex < len(flatList) {
		node := flatList[s.selectedIndex]
		return s.tree.ToggleExpanded(node)
	}
	return nil
}

// SelectedNode returns the currently selected node
func (s *Sidebar) SelectedNode() *FileNode {
	flatList := s.tree.FlatList()
	if s.selectedIndex >= 0 && s.selectedIndex < len(flatList) {
		return flatList[s.selectedIndex]
	}
	return nil
}

// Refresh the file tree while preserving the expansion state
func (s *Sidebar) Refresh() error {
	if s.tree == nil {
		return nil
	}

	// store the current state of the expanded paths
	expandedPaths := s.collectExpandedPaths(s.tree.Root)

	// Refresh the tree
	if err := s.tree.Refresh(); err != nil {
		return err
	}

	// Restore the expanded paths
	s.restoreExpandedPaths(s.tree.Root, expandedPaths)

	// Rebuild the flat list with restored state
	s.tree.rebuildFlatList()

	// Adjust selected index if its now out of bounds
	flatList := s.tree.FlatList()
	if s.selectedIndex >= len(flatList) {
		s.selectedIndex = len(flatList) - 1
	}
	if s.selectedIndex < 0 {
		s.selectedIndex = 0
	}

	s.Render()

	return nil
}

// adjustScroll adjusts the scroll offset to keep selection visible
func (s *Sidebar) adjustScroll() {
	visibleLines := s.height - 2 // Account for borders

	// Scroll down if selection is below visible area
	if s.selectedIndex >= s.scrollOffset+visibleLines {
		s.scrollOffset = s.selectedIndex - visibleLines + 1
	}

	// Scroll up if selection is above visible area
	if s.selectedIndex < s.scrollOffset {
		s.scrollOffset = s.selectedIndex
	}
}

// Recursively collects all expanded directory paths
func (s *Sidebar) collectExpandedPaths(node *FileNode) map[string]bool {
	expanded := make(map[string]bool)
	s.collectExpandedPathsRecursive(node, expanded)
	return expanded
}

// Helper to collect expanded dir paths recursively
func (s *Sidebar) collectExpandedPathsRecursive(node *FileNode, expanded map[string]bool) {
	if node.IsDir && node.Expanded {
		expanded[node.Path] = true
		for _, child := range node.Children {
			s.collectExpandedPathsRecursive(child, expanded)
		}
	}
}

// Recursively restore the expanded file path state
func (s *Sidebar) restoreExpandedPaths(node *FileNode, expandedPaths map[string]bool) {
	if !node.IsDir {
		return
	}

	// If the path was previously expanded, expand it again
	if expandedPaths[node.Path] {
		node.Expanded = true

		// Load Children if not loaded before
		if len(node.Children) == 0 {
			s.tree.ToggleExpanded(node) // load the children
			node.Expanded = true        // Ensure children are expanded
		}

		for _, child := range node.Children {
			s.restoreExpandedPaths(child, expandedPaths)
		}
	}
}
